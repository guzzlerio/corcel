package main

import (
	"fmt"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/serialisation/yaml"
	"github.com/guzzlerio/corcel/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExecutionPlanExtractions", func() {

	Context("KeyValue", func() {
		Context("Step Scope", func() {
			It("Succeeds", func() {
				var plan = fmt.Sprintf(`---
				name: Some Plan
				iterations: 0
				random: false
				workers: 1
				waitTime: 0s
				duration: 0s
				jobs:
					- name: Some Job
					  steps:
					  - name: Some Step
						action:
						  type: DummyAction
						  results:
							key: 12345
						extractors:
						 - type: KeyValueExtractor
						   key: key
						   name: target
						assertions:
						 - type: ExactAssertion
						   key: target
						   expected: 12345`)

				summary, err := test.ExecutePlanFromDataForApplication(plan)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
			})
			It("Fails", func() {

				var plan = fmt.Sprintf(`---
				name: Some Plan
				iterations: 0
				random: false
				workers: 1
				waitTime: 0s
				duration: 0s
				jobs:
				- name: Some Job
				  steps:
				  - name: Some Step
					action:
					   type: DummyAction
					   results:
						  hole: 12345
					extractors:
						- type: KeyValueExtractor
						  key: key
						  name: target
					assertions:
						- type: ExactAssertion
						  key: target
						  expected: 123456`)

				summary, err := test.ExecutePlanFromDataForApplication(plan)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})
		Context("Job Scope", func() {
			It("Succeeds", func() {
				var plan = fmt.Sprintf(`---
				name: Some Plan
				iterations: 0
				random: false
				workers: 1
				waitTime: 0s
				duration: 0s
				jobs:
					- name: Some Job
					  steps:
					  - name: Step 1
						action:
						  type: DummyAction
						  results:
							key: 12345
						extractors:
						 - type: KeyValueExtractor
						   key: key
						   name: target
						   scope: job
					  - name: Step 2
						assertions:
						 - type: ExactAssertion
						   key: target
						   expected: 12345
					  `)

				summary, err := test.ExecutePlanFromDataForApplication(plan)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertions).To(Equal(int64(1)))
				Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
			})
			It("Succeeds using extracted context variable in action", func() {
				var plan = fmt.Sprintf(`---
				name: Some Plan
				iterations: 0
				random: false
				workers: 1
				waitTime: 0s
				duration: 0s
				jobs:
					- name: Some Job
					  steps:
					  - name: Step 1
						action:
						  type: DummyAction
						  results:
							key: 12345
						extractors:
						 - type: KeyValueExtractor
						   key: key
						   name: target
						   scope: job
					  - name: Step 2
						action:
						  type: DummyAction
						  results:
							lastKey: $target
						extractors:
						 - type: KeyValueExtractor
						   key: lastKey
						   name: newKey
						   scope: job
						assertions:
						 - type: ExactAssertion
						   key: newKey
						   expected: 12345
					  `)

				summary, err := test.ExecutePlanFromDataForApplication(plan)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertions).To(Equal(int64(1)))
				Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
			})
			It("Fails", func() {

				var plan = fmt.Sprintf(`---
				name: Some Plan
				iterations: 0
				random: false
				workers: 1
				waitTime: 0s
				duration: 0s
				jobs:
					- name: Some Job
					  steps:
					  - name: Step 1
						action:
						  type: DummyAction
						  results:
							hole: 12345
						extractors:
						 - type: KeyValueExtractor
						   key: key
						   name: target
						   scope: job
					  - name:  Step 2
						assertions:
						 - type: ExactAssertion
						   key: target
						   expected: 12345
					  `)

				summary, err := test.ExecutePlanFromDataForApplication(plan)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})
		Context("Plan Scope", func() {
			It("Succeeds", func() {
				var plan = fmt.Sprintf(`---
				name: Some Plan
				iterations: 0
				random: false
				workers: 1
				waitTime: 0s
				duration: 0s
				jobs:
					- name: Job 1
					  steps:
					  - name: Step 1
						action:
						  type: DummyAction
						  results:
							key: 12345
						extractors:
						 - type: KeyValueExtractor
						   key: key
						   name: target
						   scope: plan
					- name: Job 2
					  steps: 
					  - name: Step 1
						assertions:
						 - type: ExactAssertion
						   key: target
						   expected: 12345`)

				summary, err := test.ExecutePlanFromDataForApplication(plan)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
			})
			It("Fails", func() {
				var plan = fmt.Sprintf(`---
				name: Some Plan
				iterations: 0
				random: false
				workers: 1
				waitTime: 0s
				duration: 0s
				jobs:
					- name: Job 1
					  steps:
					  - name: Step 1
						action:
						  type: DummyAction
						  results:
							hole: 12345
						extractors:
						 - type: KeyValueExtractor
						   key: key
						   name: target
						   scope: plan
					- name: Job 2
					  steps: 
					  - name: Step 1
						assertions:
						 - type: ExactAssertion
						   key: target
						   expected: 12345`)

				summary, err := test.ExecutePlanFromDataForApplication(plan)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})
	})

	Context("Regex", func() {
		Context("Step Scope", func() {
			Context("Succeeds", func() {
				It("Matches simple pattern", func() {
					planBuilder := yaml.NewPlanBuilder()

					planBuilder.
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
						WithExtractor(planBuilder.RegexExtractor().Name("regex:match:1").Key("value:1").Match("\\d+").Build()).
						WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
				})
			})
			Context("Fails", func() {
				It("Matches simple pattern", func() {
					planBuilder := yaml.NewPlanBuilder()

					planBuilder.
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
						WithExtractor(planBuilder.RegexExtractor().Name("regex:match:1").Key("value:1").Match("boom").Build()).
						WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
				})
			})

			PIt("Extends the name with any named groups", func() {})

			PIt("Extends the name with index access with any non-named groups", func() {})
		})
		Context("Job Scope", func() {
			Context("Succeeds", func() {
				It("Matches simple pattern", func() {
					planBuilder := yaml.NewPlanBuilder()

					jobBuilder := planBuilder.
						CreateJob()
					jobBuilder.
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
						WithExtractor(planBuilder.RegexExtractor().
							Name("regex:match:1").
							Key("value:1").Match("\\d+").
							Scope(core.JobScope).Build())
					jobBuilder.
						CreateStep().
						WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
				})
			})
			Context("Fails", func() {
				It("Matches simple pattern but scope not set to Job and so defaults to Step", func() {
					planBuilder := yaml.NewPlanBuilder()

					jobBuilder := planBuilder.
						CreateJob()
					jobBuilder.
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
						WithExtractor(planBuilder.RegexExtractor().
							Name("regex:match:1").
							Key("value:1").Match("\\d+").Build())
					jobBuilder.
						CreateStep().
						WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
				})
			})
		})
		Context("Plan Scope", func() {
			Context("Succeeds", func() {
				It("Matches simple pattern", func() {
					planBuilder := yaml.NewPlanBuilder()

					planBuilder.
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
						WithExtractor(planBuilder.RegexExtractor().
							Name("regex:match:1").
							Key("value:1").Match("\\d+").
							Scope(core.PlanScope).Build())

					planBuilder.
						CreateJob().
						CreateStep().
						WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
				})
			})
			Context("Fails", func() {
				It("Matches simple pattern but scope not set to Job and so defaults to Step", func() {
					planBuilder := yaml.NewPlanBuilder()

					planBuilder.
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", "talula 123 bang bang").Build()).
						WithExtractor(planBuilder.RegexExtractor().
							Name("regex:match:1").
							Key("value:1").Match("\\d+").Build())

					planBuilder.
						CreateJob().
						CreateStep().
						WithAssertion(planBuilder.ExactAssertion("regex:match:1", "123"))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
				})
			})
		})
	})

	Context("XPAth", func() {
		sampleContent := `<library>
          <!-- Great book. -->
          <book id="b0836217462" available="true">
            <isbn>0836217462</isbn>
            <title lang="en">Being a Dog Is a Full-Time Job</title>
            <quote>I'd dog paddle the deepest ocean.</quote>
            <author id="CMS">
              <?echo "go rocks"?>
              <name>Charles M Schulz</name>
              <born>1922-11-26</born>
              <dead>2000-02-12</dead>
            </author>
            <character id="PP">
              <name>Peppermint Patty</name>
              <born>1966-08-22</born>
              <qualification>bold, brash and tomboyish</qualification>
            </character>
            <character id="Snoopy">
              <name>Snoopy</name>
              <born>1950-10-04</born>
              <qualification>extroverted beagle</qualification>
            </character>
          </book>
        </library>`

		var entries = []TableEntry{
			Entry("", "/library/book/isbn", "0836217462"),
			Entry("", "library/*/isbn", "0836217462"),
			Entry("", "/library/book/../book/./isbn", "0836217462"),
			Entry("", "/library/book/character[2]/name", "Snoopy"),
			Entry("", "/library/book/character[born='1950-10-04']/name", "Snoopy"),
			Entry("", "/library/book//node()[@id='PP']/name", "Peppermint Patty"),
			Entry("", "//book[author/@id='CMS']/title", "Being a Dog Is a Full-Time Job"),
			Entry("", "/library/book/preceding::comment()", " Great book. "),
			Entry("", "//*[contains(born,'1922')]/name", "Charles M Schulz"),
			Entry("", "//*[@id='PP' or @id='Snoopy']/born", "1966-08-22"),
		}

		Context("Step Scope", func() {
			DescribeTable("Succeeds", func(testCase string, expectedValue string) {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.XPathExtractor().Name("xpath:match:1").Key("value:1").XPath(testCase).Build()).
					WithAssertion(planBuilder.ExactAssertion("xpath:match:1", expectedValue))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
			}, entries...)

			It("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.XPathExtractor().Name("xpath:match:1").Key("value:1").XPath("fubar").Build()).
					WithAssertion(planBuilder.ExactAssertion("xpath:match:1", "123"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})
		Context("Job Scope", func() {
			DescribeTable("Succeeds", func(testCase string, expectedValue string) {
				planBuilder := yaml.NewPlanBuilder()

				jobBuilder := planBuilder.
					CreateJob()
				jobBuilder.
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.XPathExtractor().Name("xpath:match:1").Key("value:1").XPath(testCase).Scope(core.JobScope).Build())
				jobBuilder.
					CreateStep().
					WithAssertion(planBuilder.ExactAssertion("xpath:match:1", expectedValue))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
			}, entries...)

			It("fails", func() {
				planBuilder := yaml.NewPlanBuilder()
				jobBuilder := planBuilder.
					CreateJob()
				jobBuilder.
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.XPathExtractor().Name("xpath:match:1").Key("value:1").XPath("/library/book/isbn").Build())
				jobBuilder.
					CreateStep().
					WithAssertion(planBuilder.ExactAssertion("xpath:match:1", "0836217462"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})

		Context("Plan Scope", func() {
			DescribeTable("Succeeds", func(testCase string, expectedValue string) {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.XPathExtractor().Name("xpath:match:1").Key("value:1").XPath(testCase).Scope(core.PlanScope).Build())

				planBuilder.
					CreateJob().
					CreateStep().
					WithAssertion(planBuilder.ExactAssertion("xpath:match:1", expectedValue))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
			}, entries...)

			It("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.XPathExtractor().Name("xpath:match:1").Key("value:1").XPath("/library/book/isbn").Build())

				planBuilder.
					CreateJob().
					CreateStep().
					WithAssertion(planBuilder.ExactAssertion("xpath:match:1", "0836217462"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})
	})

	Context("JSON Path", func() {
		var sampleContent = `{
			"store": {
				"book": [
					{
						"category": "reference",
						"author": "Nigel Rees",
						"title": "Sayings of the Century",
						"price": 8.95
					},
					{
						"category": "fiction",
						"author": "Evelyn Waugh",
						"title": "Sword of Honour",
						"price": 12.99
					},
					{
						"category": "fiction",
						"author": "Herman Melville",
						"title": "Moby Dick",
						"isbn": "0-553-21311-3",
						"price": 8.99
					},
					{
						"category": "fiction",
						"author": "J. R. R. Tolkien",
						"title": "The Lord of the Rings",
						"isbn": "0-395-19395-8",
						"price": 22.99
					}
				],
				"bicycle": {
					"color": "red",
					"price": 19.95
				}
			},
			"expensive": 10
		}`

		var entries = []TableEntry{
			Entry("", "$.expensive", float64(10)),
			Entry("", "$.store.book[0].price", float64(8.95)),
			Entry("", "$.store.book[-1].isbn", "0-395-19395-8"),
			Entry("", "$.store.book[0,1].price", []float64{8.95, 12.99}),
			Entry("", "$.store.book[0:2].price", []float64{8.95, 12.99, 8.99}),
			Entry("", "$.store.book[?(@.isbn)].price", []float64{8.99, 22.99}),
			Entry("", "$.store.book[?(@.price > 10)].title", []string{"Sword of Honour", "The Lord of the Rings"}),
			Entry("", "$.store.book[?(@.price < $.expensive)].price", []float64{8.95, 8.99}),
			Entry("", "$.store.book[:].price", []float64{8.95, 12.99, 8.99, 22.99}),
		}

		Context("Step Scope", func() {
			DescribeTable("Succeeds",
				func(testCase string, expectedValue interface{}) {
					planBuilder := yaml.NewPlanBuilder()
					planBuilder.
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
						WithExtractor(planBuilder.JSONPathExtractor().Name("jsonpath:match:1").Key("value:1").JSONPath(testCase).Build()).
						WithAssertion(planBuilder.ExactAssertion("jsonpath:match:1", expectedValue))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
				},
				entries...)

			It("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.JSONPathExtractor().Name("jsonpath:match:1").Key("value:1").JSONPath("fubar").Build()).
					WithAssertion(planBuilder.ExactAssertion("jsonpath:match:1", "123"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})
		Context("Job Scope", func() {
			DescribeTable("Succeeds",
				func(testCase string, expectedValue interface{}) {
					planBuilder := yaml.NewPlanBuilder()

					jobBuilder := planBuilder.
						CreateJob()
					jobBuilder.
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
						WithExtractor(planBuilder.JSONPathExtractor().Name("jsonpath:match:1").Key("value:1").JSONPath(testCase).Scope(core.JobScope).Build())
					jobBuilder.
						CreateStep().
						WithAssertion(planBuilder.ExactAssertion("jsonpath:match:1", expectedValue))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
				},
				entries...)

			It("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				jobBuilder := planBuilder.
					CreateJob()
				jobBuilder.
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.JSONPathExtractor().Name("jsonpath:match:1").Key("value:1").JSONPath("fubar").Build())
				jobBuilder.
					CreateStep().
					WithAssertion(planBuilder.ExactAssertion("jsonpath:match:1", "123"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})
		Context("Plan Scope", func() {
			DescribeTable("Succeeds",
				func(testCase string, expectedValue interface{}) {
					planBuilder := yaml.NewPlanBuilder()

					planBuilder.
						CreateJob().
						CreateStep().
						ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
						WithExtractor(planBuilder.JSONPathExtractor().Name("jsonpath:match:1").Key("value:1").JSONPath(testCase).Scope(core.PlanScope).Build())

					planBuilder.
						CreateJob().
						CreateStep().
						WithAssertion(planBuilder.ExactAssertion("jsonpath:match:1", expectedValue))

					summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
					Expect(err).To(BeNil())
					Expect(summary.TotalAssertionFailures).To(Equal(int64(0)))
				},
				entries...)

			It("Fails", func() {
				planBuilder := yaml.NewPlanBuilder()

				planBuilder.
					CreateJob().
					CreateStep().
					ToExecuteAction(planBuilder.DummyAction().Set("value:1", sampleContent).Build()).
					WithExtractor(planBuilder.JSONPathExtractor().Name("jsonpath:match:1").Key("value:1").JSONPath("fubar").Build())

				planBuilder.
					CreateJob().
					CreateStep().
					WithAssertion(planBuilder.ExactAssertion("jsonpath:match:1", "123"))

				summary, err := test.ExecutePlanBuilderForApplication(planBuilder)
				Expect(err).To(BeNil())
				Expect(summary.TotalAssertionFailures).To(Equal(int64(1)))
			})
		})

	})

	Context("Javascript", func() {
		PContext("Step Scope", func() {

		})
		PContext("Job Scope", func() {

		})
		PContext("Plan Scope", func() {

		})
	})
})
