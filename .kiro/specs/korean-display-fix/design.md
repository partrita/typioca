# Design Document

## Overview

The Korean display issue occurs because the current `initTimerBasedTest` function uses a generic `Generate(generatorKey)` method that doesn't properly handle Korean wordlists. When Korean wordlists are selected, the system fails to generate Korean content, resulting in blank screens during timer-based typing tests.

The solution involves modifying the content generation logic to detect Korean wordlists and use the appropriate Korean-specific generation methods that already exist in the codebase.

## Architecture

### Current Flow (Problematic)
```
User selects Korean wordlist → initTimerBasedTest() → 
mainMenu.timeBasedGenerator.Generate(generatorKey) → 
Generic Generate method → Empty/incorrect content → Blank screen
```

### Proposed Flow (Fixed)
```
User selects Korean wordlist → initTimerBasedTest() → 
Detect Korean wordlist → Use Korean-specific generator → 
GenerateKoreanWords/GenerateKoreanSentences → Korean content → Display Korean text
```

## Components and Interfaces

### 1. Korean Detection Logic
- **Function**: `isKoreanWordList(generatorKey string) bool`
- **Purpose**: Identify if the selected wordlist is Korean-based
- **Implementation**: Check if generatorKey contains "korean" prefix or Korean language indicators

### 2. Content Generation Router
- **Location**: `initTimerBasedTest` function in `cmd/init.go`
- **Purpose**: Route content generation to appropriate method based on wordlist type
- **Logic**:
  - If Korean wordlist detected → Use Korean-specific generation
  - If English wordlist detected → Use existing generic generation

### 3. Korean Content Generation
- **Existing Methods**: 
  - `GenerateKoreanWords(listName string) []rune`
  - `GenerateKoreanSentences(listName string) []rune`
- **Enhancement**: Map generatorKey to appropriate Korean generation method

### 4. Error Handling and Validation
- **Content Validation**: Verify generated content is not empty
- **Fallback Mechanism**: Provide default content if generation fails
- **Error Logging**: Log specific errors for debugging

## Data Models

### WordList Mapping
```go
type KoreanWordListMapping struct {
    GeneratorKey string
    ListName     string
    GenerationType string // "words" or "sentences"
}
```

### Content Generation Result
```go
type GenerationResult struct {
    Content []rune
    Success bool
    Error   error
}
```

## Error Handling

### 1. Empty Content Detection
- **Trigger**: When generated content is empty or nil
- **Response**: Log error and provide fallback content
- **User Experience**: Show error message instead of blank screen

### 2. Generator Initialization Failure
- **Trigger**: When Korean generator fails to initialize
- **Response**: Fall back to English generator with warning
- **User Experience**: Display warning about Korean content unavailability

### 3. Encoding Issues
- **Trigger**: When Korean characters fail to render
- **Response**: Validate UTF-8 encoding of generated content
- **User Experience**: Show encoding error message with font recommendations

## Testing Strategy

### 1. Unit Tests
- Test Korean wordlist detection logic
- Test content generation for each Korean wordlist type
- Test error handling scenarios
- Test fallback mechanisms

### 2. Integration Tests
- Test complete flow from wordlist selection to content display
- Test switching between English and Korean wordlists
- Test timer-based test initialization with Korean content

### 3. Manual Testing
- Verify Korean characters display correctly in terminal
- Test all Korean wordlist options (common words, tech terms, sentences)
- Verify timer functionality works with Korean content
- Test error scenarios (missing wordlists, encoding issues)

### 4. Edge Case Testing
- Empty Korean wordlists
- Corrupted Korean content files
- Terminal without Korean font support
- Rapid switching between wordlist types

## Implementation Approach

### Phase 1: Core Fix
1. Modify `initTimerBasedTest` to detect Korean wordlists
2. Route Korean wordlists to appropriate generation methods
3. Add content validation and error handling

### Phase 2: Enhancement
1. Improve error messages and user feedback
2. Add logging for debugging
3. Optimize Korean content generation performance

### Phase 3: Testing and Validation
1. Comprehensive testing of all Korean wordlist types
2. Validation of error handling scenarios
3. Performance testing with large Korean wordlists

## Security Considerations

- **Input Validation**: Ensure generatorKey values are sanitized
- **Content Validation**: Verify Korean content doesn't contain malicious characters
- **Error Information**: Avoid exposing sensitive system information in error messages

## Performance Considerations

- **Lazy Loading**: Load Korean wordlists only when needed
- **Caching**: Cache generated Korean content for repeated use
- **Memory Management**: Properly manage memory for large Korean sentence lists
- **Generation Speed**: Optimize Korean text generation for responsive UI