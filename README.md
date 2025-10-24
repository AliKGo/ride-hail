# ride-hail
tutorial

Transaction-manager:
    Это нужно, для того чтобы хранить согласованность данных при внесении изменении в базу данных и в rabbit
    если произойдет ошибка в каком то этапе изменение, то мы все изменение отменяем 
    как оно работает:
        Он перед запуском функции который мы передали ctx закидывает туда транзакцию. И потом этот ctx мы передаем в функцию. 
        Внутри repo мы первым делом используем метод GetExecutor из пакета executor который принимает ctx и pool *pgxpool.Pool.
        Этот метод возвращает объект реализующие интерфейс DBExecutor
        type DBExecutor interface {
            Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
            Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
            QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
        }
        это нужно, для того чтобы поддерживать вариант с транзакции и просто запуск функции отдельно без transaction-manager.

    в структуре сервиса должно быть
        txm  txm.Manager
    пример:
        type RideService struct {
            log  *logger.Logger
            repo Repository
            txm  txm.Manager
        }

    пример использвание:
	fn := func(ctx context.Context) error {
		pickupID, err := svc.repo.cord.CreateNewCoordinate(ctx, models.Coordinate{})
		if err != nil {
			return err
		}
		destinationID, err := svc.repo.cord.CreateNewCoordinate(ctx, models.Coordinate{})
		if err != nil {
			return err
		}
		rideID, err = svc.repo.ride.CreateNewRide(ctx, models.Ride{})
		if err != nil {
			return err
		}
		return nil
	}
        	err = svc.txm.Do(ctx, fn)
            fn() - это функция в которым написать сценарии цепочка сообыти который логический связан между собой
    если это функция возвращает ошибку то все изменение до ошибки будет возвращен в началное состояние