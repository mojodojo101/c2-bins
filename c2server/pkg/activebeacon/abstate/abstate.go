package abstate

import (
	"context"
	"github.com/mojodojo101/c2server/pkg/activebeacon"
	"time"
)

type activeBeaconState struct {
	activeBeaconRepo activebeacon.Repository
}

func NewActiveBeaconState(ar activebeacon.Repository) activebeacon.State {
	return &activeBeaconState{
		activeBeaconRepo: ar,
	}

}

//i messed up with psql time zone and time.time conversion so this is my ugly way to fix this locally
//what i need to do is either rewrite the entire db or some conversion stuff O:
func (as *activeBeaconState) Update(ctx context.Context, lineance int64) error {
	//amount is the amount of beacons to retrieve from the database
	amount := int64(6000)
	abs, err := as.activeBeaconRepo.GetAllActiveBeacons(ctx, amount)
	if err != nil {
		return err
	}
	for _, ab := range abs {
		t1 := time.Duration(ab.Ping)*time.Second + time.Duration(lineance)*time.Second
		t2 := time.Now().Sub(ab.UpdatedAt)
		if t2 > t1 {
			x := (time.Duration(ab.Ping) * time.Second)
			ab.MissedPings = int64(t2 / x)
		} else {
			ab.MissedPings = 0
		}
		err = as.activeBeaconRepo.Update(ctx, &ab)
		if err != nil {
			return err
		}
	}
	time.Sleep(time.Duration(lineance) * time.Second)
	return nil
}
