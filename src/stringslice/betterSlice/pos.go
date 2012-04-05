package betterSlice

type StringSlice []string

func (slice StringSlice) Pos(value string) int {
    for p, v := range slice {
        if (v == value) {
            return p
        }
    }
    return -1
}