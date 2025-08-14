package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-policy-agent/opa/v1/rego"
)

func evaluatePolicy(input map[string]any, policy string) ([]string, error) {
	ctx := context.Background()

	regoEval := rego.New(
		rego.Query("data.validate"),
		rego.Module("validate.rego", policy),
	)

	query, err := regoEval.PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	valueJSON, err := json.Marshal(results[0].Expressions[0].Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal policy result: %w", err)
	}

	var policyResult PolicyResult
	if err := json.Unmarshal(valueJSON, &policyResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal policy result: %w", err)
	}

	if len(policyResult.Violations) == 0 {
		return nil, nil
	}

	return policyResult.Violations, nil
}
