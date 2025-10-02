# Implementation Plan

- [x] 1. Create Korean word data files and basic Korean text processing
  - Create embedded Korean word list JSON files with common Korean words, tech terms, and sentences
  - Implement basic Korean character detection and validation functions
  - _Requirements: 3.1, 3.2_

- [x] 1.1 Create Korean word list JSON files
  - Create `cmd/words/embedables/words/korean-common.json` with 500+ common Korean words
  - Create `cmd/words/embedables/words/korean-tech.json` with 200+ Korean tech terms  
  - Create `cmd/words/embedables/sentences/korean-sentences.json` with 50+ Korean sentences
  - _Requirements: 3.1, 3.2_

- [x] 1.2 Implement Korean character processing utilities
  - Create `cmd/korean.go` with functions for Korean character detection (`isKoreanChar`)
  - Implement Korean WPM calculation function (`calculateKoreanWPM`) 
  - Implement Korean accuracy calculation function (`calculateKoreanAccuracy`)
  - Write unit tests for Korean character processing functions
  - _Requirements: 4.1, 4.2, 4.3_

- [x] 2. Extend words generator to support Korean content
  - Modify `cmd/words/words.go` to embed Korean word sources
  - Create Korean-specific word generator functions
  - Add Korean word source loading and processing
  - _Requirements: 1.1, 1.2, 2.1, 2.2_

- [x] 2.1 Add Korean word sources to words generator
  - Add Korean word file embedding in `cmd/words/words.go`
  - Implement `addEmbeddedKoreanSources` function to load Korean word lists
  - Modify `NewGenerator` to support Korean word sources
  - Write unit tests for Korean word source loading
  - _Requirements: 1.1, 2.1_

- [x] 2.2 Create Korean word generation functionality
  - Implement Korean word shuffling and selection logic
  - Add Korean sentence generation with proper spacing and punctuation
  - Ensure Korean text generation respects word/sentence count limits
  - Write unit tests for Korean word generation
  - _Requirements: 1.2, 2.2_

- [x] 3. Extend configuration system for Korean support
  - Add Korean word list configuration structures
  - Implement Korean word list filtering functions
  - Update config loading/saving to handle Korean settings
  - _Requirements: 5.1, 5.2, 5.3, 6.2_

- [x] 3.1 Add Korean configuration structures
  - Extend `Config` struct in `cmd/model.go` to include Korean word lists
  - Add `Language` field to `EmbededWordList` and `WordList` structs
  - Update `TestSettingCursors` to include Korean test cursors
  - _Requirements: 5.2, 6.2_

- [x] 3.2 Implement Korean word list filtering
  - Create `filterEnabledKoreanSelections`, `filterEnabledKoreanWordSelection`, `filterEnabledKoreanSentenceSelection` functions in `cmd/init.go`
  - Modify existing filter functions to distinguish between English and Korean lists
  - Write unit tests for Korean filtering functions
  - _Requirements: 5.1, 5.3_

- [x] 4. Create Korean test settings and menu options
  - Implement Korean-specific test setting structures
  - Add Korean test initialization functions
  - Update main menu to include Korean test options
  - _Requirements: 1.1, 1.2, 2.1, 2.2, 6.1, 6.3_

- [x] 4.1 Implement Korean test settings structures
  - Create `KoreanTimerBasedTestSettings`, `KoreanWordCountBasedTestSettings`, `KoreanSentenceCountBasedTestSettings` in `cmd/model.go`
  - Implement Korean test setting initialization functions in `cmd/init.go`
  - Add Korean generator fields to `MainMenu` struct
  - _Requirements: 1.1, 2.1, 6.1_

- [x] 4.2 Update main menu for Korean options
  - Modify `initMainMenu` function to include Korean test selections
  - Update menu display logic in `cmd/view.go` to show Korean options with clear language labels
  - Implement Korean test option selection handling in menu navigation
  - _Requirements: 1.1, 2.1, 6.1, 6.3_

- [x] 4.3 Implement Korean test initialization
  - Create `initKoreanTimerBasedTest`, `initKoreanWordCountBasedTest`, `initKoreanSentenceCountBasedTest` functions
  - Ensure Korean tests use Korean word generators and proper Korean text processing
  - Write unit tests for Korean test initialization
  - _Requirements: 1.2, 2.2_

- [x] 5. Implement Korean text display and input handling
  - Create Korean text wrapping and display functions
  - Implement Korean input validation and error handling
  - Update test view rendering for Korean text
  - _Requirements: 1.3, 2.3, 4.4_

- [x] 5.1 Create Korean text display functions
  - Implement `wrapKoreanStyledParagraph` function for proper Korean text wrapping
  - Create `koreanParagraphView` method for Korean text rendering
  - Update `cmd/view.go` to handle Korean text display in test views
  - _Requirements: 1.3, 2.3, 4.4_

- [x] 5.2 Implement Korean input handling
  - Add Korean character input validation in test input processing
  - Implement Korean-specific mistake tracking and highlighting
  - Update cursor positioning logic for Korean text
  - Write unit tests for Korean input handling
  - _Requirements: 1.3, 2.3_

- [x] 6. Add Korean results calculation and display
  - Implement Korean WPM and accuracy calculation in test results
  - Update results display to show Korean-specific metrics
  - Ensure Korean test results are properly formatted and saved
  - _Requirements: 4.1, 4.2, 4.3, 4.4_

- [x] 6.1 Implement Korean results calculation
  - Integrate `calculateKoreanWPM` and `calculateKoreanAccuracy` functions into test completion logic
  - Update `Results` struct creation for Korean tests to use Korean-specific calculations
  - Ensure Korean results are calculated correctly for all test types (timer, word count, sentence count)
  - _Requirements: 4.1, 4.2, 4.3_

- [x] 6.2 Update results display for Korean tests
  - Modify results view rendering in `cmd/view.go` to properly display Korean test results
  - Ensure Korean word list names and metrics are displayed correctly
  - Update results formatting to handle Korean text properly
  - _Requirements: 4.4_

- [x] 7. Add custom Korean file support
  - Implement Korean custom file loading and validation
  - Update configuration system to handle Korean custom files
  - Add Korean custom file menu integration
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [x] 7.1 Implement Korean custom file loading
  - Extend custom file loading logic in `cmd/words/words.go` to detect and process Korean files
  - Add Korean file validation to ensure proper Korean content
  - Implement Korean JSON and text file format support
  - Write unit tests for Korean custom file loading
  - _Requirements: 5.1, 5.2, 5.3_

- [x] 7.2 Update configuration for Korean custom files
  - Modify config loading/saving to handle Korean custom word lists
  - Update config view in `cmd/view.go` to display Korean custom files
  - Implement Korean custom file enable/disable functionality
  - _Requirements: 5.4_

- [x] 8. Integrate Korean tests with existing test flow
  - Update test state transitions to handle Korean tests
  - Ensure Korean tests integrate properly with existing timer, stopwatch, and results systems
  - Add Korean test restart and menu navigation functionality
  - _Requirements: 1.1, 1.2, 1.3, 2.1, 2.2, 2.3, 6.1, 6.2_

- [x] 8.1 Update test state handling for Korean tests
  - Modify test update functions in `cmd/update.go` to handle Korean test states
  - Ensure Korean tests properly integrate with timer and stopwatch systems
  - Update test completion logic to transition correctly from Korean tests to results
  - _Requirements: 1.1, 1.2, 2.1, 2.2_

- [x] 8.2 Add Korean test navigation and controls
  - Implement Korean test restart functionality (ctrl+r)
  - Ensure Korean tests can return to main menu (ctrl+q)
  - Update help text and controls display for Korean tests
  - _Requirements: 1.3, 2.3, 6.1, 6.2_

- [x] 9. Write comprehensive tests for Korean functionality
  - Create unit tests for all Korean-specific functions
  - Write integration tests for Korean typing test flows
  - Add test data and fixtures for Korean testing
  - _Requirements: 1.1, 1.2, 1.3, 2.1, 2.2, 2.3, 3.1, 3.2, 4.1, 4.2, 4.3, 4.4, 5.1, 5.2, 5.3, 5.4, 6.1, 6.2, 6.3_

- [x] 9.1 Create Korean unit tests
  - Write tests for Korean character processing functions in `cmd/korean_test.go`
  - Create tests for Korean word generation and filtering
  - Add tests for Korean configuration handling
  - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1_

- [x] 9.2 Write Korean integration tests
  - Create end-to-end tests for Korean typing test scenarios
  - Test Korean custom file loading and usage
  - Verify Korean results calculation and display
  - _Requirements: 1.2, 1.3, 2.2, 2.3, 4.2, 4.3, 4.4, 5.2, 5.3, 5.4, 6.1, 6.2, 6.3_