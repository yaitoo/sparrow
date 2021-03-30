# fmt

Package fmt implements a simple fmt wrapper with higher performance.
It will hit performance issue when `fmt.Sprintf`,`fmt.Printf` or `fmt.Fprintf` is caled too many times.Because it works based on `reflect`.
We wrap standard `fmt`, and improves performance by caching and resuing reflect result.