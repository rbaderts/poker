package domain

import "github.com/jackc/pgx/v4/pgxpool"

var GblUserRepository UserRepository
var GblAccountRepository AccountRepository
var GblTableRepository TableRepository
var GblMessageDispatcher MessageDispatcher
var GblGameStateRepository GameStateRepository
var GblTransactionRepository TransactionRepository
var GblRecoveryRepository RecoveryRepository
var GblDBPool *pgxpool.Pool

