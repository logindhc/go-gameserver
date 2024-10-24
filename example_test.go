package main

import (
	"fmt"
	"gameserver/common/persistence/cache"
	"gameserver/common/utils"
	"gameserver/models"
	"sync"
	"testing"
	"time"

	list "github.com/duke-git/lancet/v2/datastructure/list"
)

func TestBuffer(t *testing.T) {
	repository := models.AccountRepository
	account := models.Account{}
	account = models.Account{ID: "1", Channel: 1, OpenId: ""}
	a := repository.GetOrCreate(account.ID)
	fmt.Printf("add get %+v \n", a)
	account.TotalLoginDay = account.TotalLoginDay + 1
	repository.Update(&account)
	repository.Flush()
	a = repository.GetOrCreate(account.ID)
	fmt.Printf("update get %+v \n", a)
	//repository.Remove(account.ID)
	//a = repository.Get(account.ID)
	//fmt.Printf("remove get %+v \n", a)
	all := repository.GetAll()
	for _, a2 := range all {
		fmt.Printf("%+v \n", a2)
	}
	//startTime := time.Now().Unix()
	//var wg sync.WaitGroup
	//for i := 0; i < 10; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		a := repository.GetOrCreate("dhc12234")
	//		//a := repository.GetOrCreate(fmt.Sprintf("dhc1_%d", i))
	//		a.OpenId = fmt.Sprintf("dhc1_%d", i)
	//		a.TotalLoginDay = a.TotalLoginDay + 1
	//		repository.Update(a)
	//		fmt.Printf("%d -%p %+v \n", i, a, a)
	//		//fmt.Printf("%d --%p %+v \n", i, a, a)
	//	}()
	//}
	//wg.Wait()
	//repository.Flush()
	//end := time.Now().Unix()
	//fmt.Printf("cost %d %d \n", end, startTime)
}

func TestCache(t *testing.T) {
	lruCache := cache.NewLRUCache[string, models.Account](200*time.Millisecond, 10*time.Millisecond)
	account := models.Account{ID: "1", Channel: 1, OpenId: ""}
	lruCache.Put(account.ID, &account)
	fmt.Printf("put %+v size %d \n", lruCache.Get(account.ID), lruCache.Size())
	//time.Sleep(110 * time.Millisecond)
	lruCache.Remove(account.ID)
	fmt.Printf("remove %+v size %d \n", lruCache.Get(account.ID), lruCache.Size())
	lruCache.Clear()
	lruCache.Flush()
}

func TestLogger(t *testing.T) {
	repository := log.DotLoginRepository
	unix := time.Now().Unix()
	//Int := new(atomic.Int32)
	total := 1
	Int := 0
	//s := sync.Mutex{}
	var wg sync.WaitGroup
	//for i := 0; i < 2000; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		//Int.Add(1)
	//		s.Lock()
	//		defer s.Unlock()
	//		Int++
	//		repository.Add(&log.DotLogin{ID: 2000, DayIndex: utils.GetYYYYMMDD(), FirstTime: &unix, LastTime: &unix, TotalCount: &total})
	//	}()
	//}
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//s.Lock()
			//defer s.Unlock()
			Int++
			//Int.Add(1)
			repository.Add(&log.DotLogin{ID: 2000, DayIndex: utils.GetYYYYMMDD(), FirstTime: &unix, LastTime: &unix, TotalCount: &total})
		}()
	}
	//
	wg.Wait()
	//wg.Add(1000)
	//for i := 0; i < 1000; i++ {
	//go func() {
	//	defer wg.Done()
	//repository.Add(&log.DotLogin{ID: 1111, DayIndex: utils.GetYYYYMMDD(), FirstTime: &unix, LastTime: &unix, TotalCount: &total})
	//}()
	//Int++
	//}
	//wg.Wait()
	fmt.Printf("cost %v \n", Int)
}

func TestUUID(t *testing.T) {
	list := list.NewList([]int{})
	//var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		//wg.Add(1)
		//go func() {
		//defer wg.Done()
		list.Push(i)
		//}()
	}
	//wg.Wait()
	newList := list.Data()
	for i, l := range newList {
		fmt.Printf("%d - %d \n", i, l)
	}
	list.Clear()
	newList2 := list.Data()
	for i, l := range newList2 {
		fmt.Printf("--%d - %d \n", i, l)
	}
}

func TestConcurrentSet(t *testing.T) {
	repository := log.DotDeviceRepository
	orCreate := repository.GetOrCreate("ae423424")
	country := int64(1)
	orCreate.Country = &country
	repository.Update(orCreate)

}
