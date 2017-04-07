package main

import (
	"testing"

	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExecutionPlan_Assertions(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("ExecutionPlan Assertions", t, func() {

		Convey("ExactAssertion", func() {

			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula").Build()).
					WithAssertion(planBuilder.ExactAssertion("value:1", "talula"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
					WithAssertion(planBuilder.ExactAssertion("value:1", 1))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

		})

		Convey("EmptyAssertion", func() {

			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", "").Build()).
					WithAssertion(planBuilder.EmptyAssertion("value:1"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", "1").Build()).
					WithAssertion(planBuilder.EmptyAssertion("value:1"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

		})

		Convey("GreaterThanAssertion", func() {

			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
					WithAssertion(planBuilder.GreaterThanAssertion("value:1", 2))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
					WithAssertion(planBuilder.GreaterThanAssertion("value:1", 5))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

		})

		Convey("GreaterThanOrEqualAssertion", func() {

			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
					WithAssertion(planBuilder.GreaterThanOrEqualAssertion("value:1", 5))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 2).Build()).
					WithAssertion(planBuilder.GreaterThanOrEqualAssertion("value:1", 5))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

		})

		Convey("LessThanAssertion", func() {

			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 3).Build()).
					WithAssertion(planBuilder.LessThanAssertion("value:1", 5))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
					WithAssertion(planBuilder.LessThanAssertion("value:1", 3))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

		})

		Convey("LessThanOrEqualAssertion", func() {

			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
					WithAssertion(planBuilder.LessThanOrEqualAssertion("value:1", 5))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
					WithAssertion(planBuilder.LessThanOrEqualAssertion("value:1", 4))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

		})

		Convey("NotEmptyAssertion", func() {

			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
					WithAssertion(planBuilder.NotEmptyAssertion("value:1"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:2", 5).Build()).
					WithAssertion(planBuilder.NotEmptyAssertion("value:1"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

		})

		Convey("NotEqualAssertion", func() {

			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 5).Build()).
					WithAssertion(planBuilder.NotEqualAssertion("value:1", 6))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", 6).Build()).
					WithAssertion(planBuilder.NotEqualAssertion("value:1", 6))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(1))
			})

		})
	})
}
