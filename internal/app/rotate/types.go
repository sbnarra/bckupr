package rotate

type Period struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type RotatePolicy struct {
	Period Period `json:"period"`
	Keep   int    `json:"keep"`
}

type RotatePolicies struct {
	Policies []RotatePolicy `json:"policies"`
}
