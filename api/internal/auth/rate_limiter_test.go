package auth

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RateLimiterTestSuite defines the test suite for RateLimiter
type RateLimiterTestSuite struct {
	suite.Suite
	rateLimiter *RateLimiter
}

// SetupTest runs before each test
func (suite *RateLimiterTestSuite) SetupTest() {
	// Create rate limiter with 3 attempts per 10 seconds, block for 5 seconds
	suite.rateLimiter = NewRateLimiter(3, 10*time.Second, 5*time.Second)
}

// TearDownTest runs after each test
func (suite *RateLimiterTestSuite) TearDownTest() {
	// Clean up any resources if needed
}

// Test NewRateLimiter
func (suite *RateLimiterTestSuite) TestNewRateLimiter() {
	rl := NewRateLimiter(5, time.Minute, time.Hour)
	
	assert.NotNil(suite.T(), rl)
	assert.Equal(suite.T(), 5, rl.maxAttempts)
	assert.Equal(suite.T(), time.Minute, rl.window)
	assert.Equal(suite.T(), time.Hour, rl.blockTime)
	assert.NotNil(suite.T(), rl.attempts)
}

// Test IsAllowed with new identifier
func (suite *RateLimiterTestSuite) TestIsAllowed_NewIdentifier() {
	identifier := "test@example.com"
	
	allowed := suite.rateLimiter.IsAllowed(identifier)
	
	assert.True(suite.T(), allowed)
}

// Test IsAllowed within limits
func (suite *RateLimiterTestSuite) TestIsAllowed_WithinLimits() {
	identifier := "test@example.com"
	
	// Record 2 attempts (under limit of 3)
	suite.rateLimiter.RecordAttempt(identifier)
	suite.rateLimiter.RecordAttempt(identifier)
	
	allowed := suite.rateLimiter.IsAllowed(identifier)
	
	assert.True(suite.T(), allowed)
}

// Test IsAllowed at limit
func (suite *RateLimiterTestSuite) TestIsAllowed_AtLimit() {
	identifier := "test@example.com"
	
	// Record 3 attempts (at limit)
	suite.rateLimiter.RecordAttempt(identifier)
	suite.rateLimiter.RecordAttempt(identifier)
	suite.rateLimiter.RecordAttempt(identifier)
	
	// Should be blocked now
	allowed := suite.rateLimiter.IsAllowed(identifier)
	
	assert.False(suite.T(), allowed)
}

// Test IsAllowed over limit
func (suite *RateLimiterTestSuite) TestIsAllowed_OverLimit() {
	identifier := "test@example.com"
	
	// Record 4 attempts (over limit of 3)
	for i := 0; i < 4; i++ {
		suite.rateLimiter.RecordAttempt(identifier)
	}
	
	allowed := suite.rateLimiter.IsAllowed(identifier)
	
	assert.False(suite.T(), allowed)
}

// Test RecordAttempt first attempt
func (suite *RateLimiterTestSuite) TestRecordAttempt_FirstAttempt() {
	identifier := "test@example.com"
	
	suite.rateLimiter.RecordAttempt(identifier)
	
	suite.rateLimiter.mutex.RLock()
	record, exists := suite.rateLimiter.attempts[identifier]
	suite.rateLimiter.mutex.RUnlock()
	
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), 1, record.count)
	assert.False(suite.T(), record.firstTry.IsZero())
	assert.True(suite.T(), record.blockedAt.IsZero())
}

// Test RecordAttempt multiple attempts
func (suite *RateLimiterTestSuite) TestRecordAttempt_MultipleAttempts() {
	identifier := "test@example.com"
	
	suite.rateLimiter.RecordAttempt(identifier)
	suite.rateLimiter.RecordAttempt(identifier)
	
	suite.rateLimiter.mutex.RLock()
	record := suite.rateLimiter.attempts[identifier]
	suite.rateLimiter.mutex.RUnlock()
	
	assert.Equal(suite.T(), 2, record.count)
	assert.True(suite.T(), record.blockedAt.IsZero()) // Not blocked yet
}

// Test RecordAttempt blocking
func (suite *RateLimiterTestSuite) TestRecordAttempt_Blocking() {
	identifier := "test@example.com"
	
	// Record attempts up to limit
	for i := 0; i < 3; i++ {
		suite.rateLimiter.RecordAttempt(identifier)
	}
	
	suite.rateLimiter.mutex.RLock()
	record := suite.rateLimiter.attempts[identifier]
	suite.rateLimiter.mutex.RUnlock()
	
	assert.Equal(suite.T(), 3, record.count)
	assert.False(suite.T(), record.blockedAt.IsZero()) // Should be blocked now
}

// Test Reset
func (suite *RateLimiterTestSuite) TestReset() {
	identifier := "test@example.com"
	
	// Record some attempts
	suite.rateLimiter.RecordAttempt(identifier)
	suite.rateLimiter.RecordAttempt(identifier)
	
	// Verify attempts exist
	suite.rateLimiter.mutex.RLock()
	_, exists := suite.rateLimiter.attempts[identifier]
	suite.rateLimiter.mutex.RUnlock()
	assert.True(suite.T(), exists)
	
	// Reset
	suite.rateLimiter.Reset(identifier)
	
	// Verify attempts are cleared
	suite.rateLimiter.mutex.RLock()
	_, exists = suite.rateLimiter.attempts[identifier]
	suite.rateLimiter.mutex.RUnlock()
	assert.False(suite.T(), exists)
}

// Test window expiration
func (suite *RateLimiterTestSuite) TestWindowExpiration() {
	// Create rate limiter with very short window
	rl := NewRateLimiter(1, 100*time.Millisecond, time.Second)
	identifier := "test@example.com"
	
	// Record attempt to reach limit
	rl.RecordAttempt(identifier)
	
	// Should be at limit
	assert.False(suite.T(), rl.IsAllowed(identifier))
	
	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)
	
	// Should be allowed again
	assert.True(suite.T(), rl.IsAllowed(identifier))
}

// Test block time expiration (basic test - timing sensitive tests are covered by StateAfterBlockTime)
func (suite *RateLimiterTestSuite) TestBlockTimeExpiration() {
	// Create rate limiter with short block time
	rl := NewRateLimiter(1, time.Second, 50*time.Millisecond)
	identifier := "test@example.com"
	
	// Exceed limit to get blocked
	rl.RecordAttempt(identifier)
	
	// Should be blocked
	assert.False(suite.T(), rl.IsAllowed(identifier))
	
	// This test mainly verifies the blocking mechanism works
	// Detailed timing tests are covered in TestStateAfterBlockTime
}

// Test concurrent access
func (suite *RateLimiterTestSuite) TestConcurrentAccess() {
	identifier := "test@example.com"
	numGoroutines := 10
	attemptsPerGoroutine := 5
	
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	
	// Launch multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < attemptsPerGoroutine; j++ {
				suite.rateLimiter.IsAllowed(identifier)
				suite.rateLimiter.RecordAttempt(identifier)
				time.Sleep(time.Millisecond) // Small delay
			}
		}()
	}
	
	wg.Wait()
	
	// Verify state is consistent
	suite.rateLimiter.mutex.RLock()
	record, exists := suite.rateLimiter.attempts[identifier]
	suite.rateLimiter.mutex.RUnlock()
	
	assert.True(suite.T(), exists)
	assert.Greater(suite.T(), record.count, 0)
}

// Test multiple identifiers
func (suite *RateLimiterTestSuite) TestMultipleIdentifiers() {
	identifier1 := "user1@example.com"
	identifier2 := "user2@example.com"
	
	// Record attempts for first identifier
	suite.rateLimiter.RecordAttempt(identifier1)
	suite.rateLimiter.RecordAttempt(identifier1)
	suite.rateLimiter.RecordAttempt(identifier1)
	
	// First identifier should be blocked
	assert.False(suite.T(), suite.rateLimiter.IsAllowed(identifier1))
	
	// Second identifier should still be allowed
	assert.True(suite.T(), suite.rateLimiter.IsAllowed(identifier2))
	
	// Record attempts for second identifier
	suite.rateLimiter.RecordAttempt(identifier2)
	
	// Second identifier should still be allowed (only 1 attempt)
	assert.True(suite.T(), suite.rateLimiter.IsAllowed(identifier2))
}

// Test window reset after expiration
func (suite *RateLimiterTestSuite) TestWindowResetAfterExpiration() {
	// Create rate limiter with short window
	rl := NewRateLimiter(2, 100*time.Millisecond, time.Second)
	identifier := "test@example.com"
	
	// Record attempts to fill window
	rl.RecordAttempt(identifier)
	rl.RecordAttempt(identifier)
	
	// Should be at limit
	assert.False(suite.T(), rl.IsAllowed(identifier))
	
	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)
	
	// Record new attempt (should reset window)
	rl.RecordAttempt(identifier)
	
	// Check internal state
	rl.mutex.RLock()
	record := rl.attempts[identifier]
	rl.mutex.RUnlock()
	
	// Count should be reset to 1
	assert.Equal(suite.T(), 1, record.count)
	assert.True(suite.T(), record.blockedAt.IsZero())
}

// Test edge case: zero attempts
func (suite *RateLimiterTestSuite) TestZeroAttempts() {
	rl := NewRateLimiter(0, time.Minute, time.Minute)
	identifier := "test@example.com"
	
	// Should be allowed initially (no attempts recorded yet)
	assert.True(suite.T(), rl.IsAllowed(identifier))
	
	// Recording attempt should block immediately
	rl.RecordAttempt(identifier)
	assert.False(suite.T(), rl.IsAllowed(identifier))
}

// Test cleanup functionality (basic test)
func (suite *RateLimiterTestSuite) TestCleanupExists() {
	// This is a basic test to ensure cleanup goroutine starts
	// Full cleanup testing would require more complex time manipulation
	
	rl := NewRateLimiter(3, time.Second, time.Second)
	
	// Add some attempts
	rl.RecordAttempt("test1@example.com")
	rl.RecordAttempt("test2@example.com")
	
	// Verify attempts exist
	rl.mutex.RLock()
	count := len(rl.attempts)
	rl.mutex.RUnlock()
	
	assert.Equal(suite.T(), 2, count)
}

// Test successful login resets attempts
func (suite *RateLimiterTestSuite) TestSuccessfulLoginReset() {
	identifier := "test@example.com"
	
	// Record failed attempts to reach limit
	suite.rateLimiter.RecordAttempt(identifier)
	suite.rateLimiter.RecordAttempt(identifier)
	suite.rateLimiter.RecordAttempt(identifier)
	
	// Verify attempts blocked
	assert.False(suite.T(), suite.rateLimiter.IsAllowed(identifier))
	
	// Simulate successful login
	suite.rateLimiter.Reset(identifier)
	
	// Should be allowed again
	assert.True(suite.T(), suite.rateLimiter.IsAllowed(identifier))
}

// Test rate limiter state after block time
func (suite *RateLimiterTestSuite) TestStateAfterBlockTime() {
	// Create rate limiter with short times for testing
	rl := NewRateLimiter(1, 50*time.Millisecond, 100*time.Millisecond)
	identifier := "test@example.com"
	
	// Exceed limit
	rl.RecordAttempt(identifier)
	assert.False(suite.T(), rl.IsAllowed(identifier))
	
	// Wait for block time to expire
	time.Sleep(150 * time.Millisecond)
	
	// Should be allowed again
	assert.True(suite.T(), rl.IsAllowed(identifier))
	
	// But window should have also expired, so count should reset
	rl.RecordAttempt(identifier)
	
	rl.mutex.RLock()
	record := rl.attempts[identifier]
	rl.mutex.RUnlock()
	
	assert.Equal(suite.T(), 1, record.count)
}

// Run the test suite
func TestRateLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}

// Additional unit tests for edge cases
func TestRateLimiter_EdgeCases(t *testing.T) {
	// Test with very large values
	rl := NewRateLimiter(1000000, 24*time.Hour, 24*time.Hour)
	assert.NotNil(t, rl)
	
	// Test with zero block time
	rl2 := NewRateLimiter(3, time.Minute, 0)
	assert.NotNil(t, rl2)
	
	identifier := "test@example.com"
	rl2.RecordAttempt(identifier)
	rl2.RecordAttempt(identifier)
	rl2.RecordAttempt(identifier)
	
	// Should be blocked but with zero block time
	assert.False(t, rl2.IsAllowed(identifier))
}

func TestRateLimiter_MemoryUsage(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute, time.Minute)
	
	// Add many identifiers
	for i := 0; i < 1000; i++ {
		identifier := fmt.Sprintf("user%d@example.com", i)
		rl.RecordAttempt(identifier)
	}
	
	rl.mutex.RLock()
	count := len(rl.attempts)
	rl.mutex.RUnlock()
	
	assert.Equal(t, 1000, count)
}