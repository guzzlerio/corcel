package main

import (
    "time"
)

func Time(function func()) time.Duration {
    now := time.Now()
    function()
    return time.Since(now)
}

func DurationIsBetween(actual time.Duration, min time.Duration, max time.Duration) bool {
    return actual >= min && actual < max
}
