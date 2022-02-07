package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Describe("normalizeDescription()", func() {
		DescribeTable("normalization of markdown",
			func(input, output string) {
				Expect(normalizeDescription(input)).To(Equal(output))
			},
			Entry("Blank input", "", ""),
			Entry("Plain text input", "Hello world", "Hello world"),
			Entry("Markdown input without comments", "**Hello world**\n\n* test\n*test2", "**Hello world**\n\n* test\n*test2"),
			Entry("Markdown with comments", "**Hello world**\n\n<!--- remove this if no breaking changes -->", "**Hello world**"),
		)
	})
})
