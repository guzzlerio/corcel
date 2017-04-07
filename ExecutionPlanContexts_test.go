package main

import (
	"os"
	"testing"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExecutionPlanContexts(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("ExecutionPlanContexts", t, func() {

		contextsDebugPath := "/tmp/contexts"
		Convey("Lists", func() {

			func() {
				if err := os.Remove(contextsDebugPath); err != nil {

				}
			}()

			defer func() {
				if err := os.Remove(contextsDebugPath); err != nil {

				}
			}()

			Convey("Plan Scope", func() {

				Convey("Succeeds", func() {
					planBuilder := yaml.NewPlanBuilder()

					planBuilder.
						SetIterations(3).
						WithContext(planBuilder.BuildContext().SetList("People", []core.ExecutionContext{
							{"name": "jill", "age": 35},
							{"name": "bob", "age": 52},
							{"name": "carol", "age": 24},
						}).Build()).
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().LogToFile(contextsDebugPath).Build())

					_, err := test.ExecutePlanBuilderForApplication(planBuilder)
					So(err, ShouldBeNil)

					contexts := test.GetExecutionContexts(contextsDebugPath)
					So(len(contexts), ShouldEqual, 3)
					So(contexts[0]["$People.name"], ShouldEqual, "jill")
					So(contexts[1]["$People.name"], ShouldEqual, "bob")
					So(contexts[2]["$People.name"], ShouldEqual, "carol")
				})

			})
		})

		Convey("Plan Scope", func() {
			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				jobBuilder := planBuilder.
					WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build()).
					CreateJob()

				jobBuilder.CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build())

				jobBuilder.CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build()).
					WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
					WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
					WithAssertion(planBuilder.ExactAssertion("$value:3", "3")).
					WithAssertion(planBuilder.ExactAssertion("something", "1"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})
			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Build()).
					WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
					WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
					WithAssertion(planBuilder.ExactAssertion("$value:3", "3"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(3))
			})
		})

		Convey("Job Scope", func() {
			Convey("Succeeds", func() {
				planBuilder := yaml.NewPlanBuilder()

				jobBuilder := planBuilder.
					CreateJob().
					WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build())

				jobBuilder.CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build())

				jobBuilder.CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("something", "$value:1").Build()).
					WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
					WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
					WithAssertion(planBuilder.ExactAssertion("$value:3", "3")).
					WithAssertion(planBuilder.ExactAssertion("something", "1"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(0))
			})

			Convey("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					WithContext(planBuilder.BuildContext().Set("value:1", "1").Set("value:2", "2").Set("value:3", "3").Build())

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Build()).
					WithAssertion(planBuilder.ExactAssertion("$value:1", "1")).
					WithAssertion(planBuilder.ExactAssertion("$value:2", "2")).
					WithAssertion(planBuilder.ExactAssertion("$value:3", "3"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				So(err, ShouldBeNil)
				So(summary.TotalAssertionFailures, ShouldEqual, int64(3))
			})
		})

	})
}
