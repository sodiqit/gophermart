package accrual

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/repository"
)

type AccrualProcessor interface {
	Run(ctx context.Context) error
}

type OrderProcessor struct {
	poolSize   int
	orderQueue chan string
	orderRepo  repository.OrderRepository
	wg         sync.WaitGroup
	logger     logger.Logger
	client     AccrualClient
}

func (p *OrderProcessor) worker(ctx context.Context, workerID int) {
	defer p.wg.Done()

	logger := p.logger.With("workerID", workerID)

	for {
		select {
		case <-ctx.Done():
			return
		case orderID, ok := <-p.orderQueue:
			if !ok {
				return
			}
			logger.Debugw("process order", "orderID", orderID)

			result, err := p.client.GetOrderInfo(ctx, orderID)

			if err != nil {
				if !(errors.Is(err, ErrOrderNotFound) || errors.Is(err, ErrRateLimit)) {
					logger.Errorw("failed to get order info", "err", err)
				}
				p.wg.Done()
				continue
			}

			err = p.orderRepo.UpdateOrder(ctx, result.OrderID, result.Status, result.Accrual)

			if err != nil {
				logger.Errorw("failed to update order", "err", err)
			} else {
				logger.Debugw("success process order", "orderID", orderID)
			}

			p.wg.Done()
		}
	}
}

func (p *OrderProcessor) Run(ctx context.Context) error {
	for i := 0; i < p.poolSize; i++ {
		go p.worker(ctx, i)
	}

	p.logger.Infow("start processing orders")

	for {
		select {
		case <-ctx.Done():
			p.wg.Wait()
			close(p.orderQueue)
			p.logger.Infow("stop processing orders")
			return ctx.Err()
		default:
			orderList, err := p.orderRepo.GetOrdersForProcessing(ctx, int64(p.poolSize)) // TODO: handle if orders deadlock in accrual system

			if err != nil || len(orderList) == 0 {
				time.Sleep(5 * time.Second)
				continue
			}

			p.wg.Add(len(orderList))
			for _, orderID := range orderList {
				p.orderQueue <- orderID
			}
			p.wg.Wait()
		}
	}
}

func NewOrderProcessor(poolSize int, orderRepo repository.OrderRepository, logger logger.Logger, client AccrualClient) *OrderProcessor {
	return &OrderProcessor{
		poolSize:   poolSize,
		orderRepo:  orderRepo,
		orderQueue: make(chan string, poolSize),
		wg:         sync.WaitGroup{},
		logger:     logger,
		client:     client,
	}
}
