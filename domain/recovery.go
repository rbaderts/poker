package domain

import (
	"context"
	"fmt"
)

type RecoveryRepository interface {

	CreateTableTombstone(ctx context.Context, tableId int) error
	DeleteTableTombstone(ctx context.Context, tableId int) error
	ReadTableTombstones(ctx context.Context) ([]int, error)

}


/**

return s a map of userId to amount recovered
 */
func RecoverTables() map[int]int {
	ctx := context.Background()
	tables, err := GblRecoveryRepository.ReadTableTombstones(ctx)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	for _, tableId := range tables {
		GblRecoveryRepository.DeleteTableTombstone(ctx, tableId)
	}

    return nil
}

func  RecoverTable(tableId int) map[int]int {

	fmt.Printf("Recovering table %d\n", tableId)
	ctx := context.Background()
	playerStates, err := GblTableRepository.GetAllPlayerStates(ctx, tableId)


	res := make(map[int]int,0)
	for _, ps := range playerStates {
		userId := ps.UserId
		delta := ps.StackChange
		starting := ps.StartingStack
		returnToPlayer := starting + delta
		fmt.Printf("return to player %d, $d dollars\n", ps.UserId, returnToPlayer)
		if returnToPlayer > 0 {
			err = GblUserRepository.AddMoney(ctx, userId, returnToPlayer)
			res[userId] += returnToPlayer

		} else if returnToPlayer <= 0 {
			fmt.Printf("returnToPlayer is negative money\n")
		}
		if err != nil {
			fmt.Errorf("error recovery funds for user %d from table %d\n", userId, tableId)
		}

	}

    GblTableRepository.ClearPlayerStates(ctx, tableId)

    return res
}