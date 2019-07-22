package main

func removeDuplicate(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	j := 0
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		in[j] = s
		j++
	}
	return in[:j]
}
