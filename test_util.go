// The MIT License (MIT)
// Copyright (C) 2019 Georgy Komarov <jubnzv@gmail.com>

package tmux

import ("os")

func InTravis() bool {
	if os.Getenv("IN_TRAVIS") == "1" {
        return true
    } else {
        return false
    }
}
