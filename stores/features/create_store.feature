Feature: Create Store

  As a store owner
  I should be able to create a store

	Scenario: Create a store called "Walmart"
		Given a valid store owner
    And no store called "Walmart" exists
    When I create a store called "Walmart"
    Then a store called "Walmart" exists
