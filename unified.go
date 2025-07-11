package tconv

// This file re-exports the public types from the models package
// to provide a clean API for library consumers.

import (
	"github.com/tpyle/tconv/pkg/models"
)

// Public type aliases for the unified test result format.
// These are the main data structures that external libraries should use.

type (
	// UnifiedTestResult represents the standardized test result format.
	UnifiedTestResult = models.UnifiedTestResult
	
	// TestMetadata contains contextual information about the test run.
	TestMetadata = models.TestMetadata
	
	// TestSuite represents a logical grouping of related tests.
	TestSuite = models.TestSuite
	
	// TestCase represents an individual test execution.
	TestCase = models.TestCase
	
	// TestStatus represents the possible outcomes of a test execution.
	TestStatus = models.TestStatus
	
	// Assertion represents an individual assertion within a test case.
	Assertion = models.Assertion
	
	// TestSummary provides aggregate statistics for a test result.
	TestSummary = models.TestSummary
)

// Status constants for test outcomes.
const (
	StatusPassed  = models.StatusPassed  // Test completed successfully
	StatusFailed  = models.StatusFailed  // Test failed assertion(s)
	StatusSkipped = models.StatusSkipped // Test was not executed  
	StatusError   = models.StatusError   // Test encountered a runtime error
)