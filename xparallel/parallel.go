package xparallel

import "sync"

func worker[V any](wg *sync.WaitGroup, ch chan V, fn func(V)) {
	for v := range ch {
		fn(v)
	}
	wg.Done()
}

func closeThenParallel[V any](maxp int, ch chan V, fn func(V)) {
	close(ch)
	concurrency := len(ch)
	if concurrency > maxp {
		concurrency = maxp
	}
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 1; i < concurrency; i++ {
		go worker(&wg, ch, fn)
	}
	worker(&wg, ch, fn)
	wg.Wait()
}

func MapK[K comparable, V any](maxp int, p map[K]V, fn func(K)) {
	ch := make(chan K, len(p))
	for k := range p {
		ch <- k
	}
	closeThenParallel(maxp, ch, fn)
}

func MapV[K comparable, V any](maxp int, p map[K]V, fn func(V)) {
	ch := make(chan V, len(p))
	for _, v := range p {
		ch <- v
	}
	closeThenParallel(maxp, ch, fn)
}

func SliceK[V any](maxp int, p []V, fn func(int)) {
	ch := make(chan int, len(p))
	for k := range p {
		ch <- k
	}
	closeThenParallel(maxp, ch, fn)
}

func SliceV[V any](maxp int, p []V, fn func(V)) {
	ch := make(chan V, len(p))
	for _, v := range p {
		ch <- v
	}
	closeThenParallel(maxp, ch, fn)
}
