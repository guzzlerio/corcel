package inproc_test

import (
	"context"

	"github.com/guzzlerio/corcel/core"
	. "github.com/guzzlerio/corcel/infrastructure/inproc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DummyAction", func() {

	PIt("Execute with string property", func() {
		var action = DummyAction{
			Results: map[string]interface{}{
				"Data": "$value",
			},
		}

		var executionContext = core.ExecutionContext{
			"vars": core.ExecutionContext{
				"$value": "talula",
			},
		}

		var result = action.Execute(context.TODO(), executionContext)

		Expect(result["Data"]).To(Equal("talula"))
	})

	PIt("Execute with float64 poerty", func() {
		var action = DummyAction{
			Results: map[string]interface{}{
				"Data": "$value",
			},
		}

		var executionContext = core.ExecutionContext{
			"vars": core.ExecutionContext{
				"$value": float64(100),
			},
		}

		var result = action.Execute(context.TODO(), executionContext)

		Expect(result["Data"]).To(Equal("100"))

	})
})
