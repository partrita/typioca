# Design Document

## Overview

이 설계는 기존 typioca 타이핑 테스트 도구에 한글 지원을 추가하는 기능을 정의합니다. 기존 아키텍처를 최대한 활용하면서 한글 문자의 특성(조합형 문자, 다양한 자모 조합)을 고려한 타이핑 테스트 환경을 제공합니다.

기존 시스템은 영어 단어와 문장을 기반으로 한 타이핑 테스트를 제공하며, 단어 생성기(WordsGenerator), 설정 관리(Config), 테스트 상태 관리(TestBase) 등의 컴포넌트로 구성되어 있습니다. 한글 지원은 이러한 기존 구조를 확장하여 구현됩니다.

## Architecture

### 기존 아키텍처 분석

현재 시스템의 주요 컴포넌트:
- `words.WordsGenerator`: 단어/문장 생성 및 관리
- `Config`: 설정 및 단어 목록 관리  
- `MainMenu`: 메뉴 시스템 및 테스트 모드 선택
- `TestBase`: 타이핑 테스트 핵심 로직
- `WordsSelection`: 단어 목록 선택 구조체

### 한글 지원 확장 아키텍처

```
┌─────────────────────────────────────────────────────────────┐
│                    Main Menu System                          │
├─────────────────────────────────────────────────────────────┤
│  English Tests          │         Korean Tests               │
│  ├─Timer Based          │         ├─Timer Based (한글)       │
│  ├─Word Count           │         ├─Word Count (한글)        │
│  └─Sentence Count       │         └─Sentence Count (한글)    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                Korean Words Generator                        │
├─────────────────────────────────────────────────────────────┤
│  ├─Embedded Korean Word Lists                               │
│  │  ├─Common Korean Words (일반 한글 단어)                   │
│  │  ├─Korean Sentences (한글 문장)                          │
│  │  └─Korean Tech Terms (한글 기술 용어)                    │
│  └─Custom Korean Word Lists                                 │
│     ├─JSON Format Support                                   │
│     └─Plain Text Format Support                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                Korean Text Processing                        │
├─────────────────────────────────────────────────────────────┤
│  ├─Korean Character Input Handling                          │
│  ├─Korean WPM Calculation                                   │
│  ├─Korean Accuracy Measurement                              │
│  └─Korean Text Display & Wrapping                          │
└─────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. Korean Word Lists (한글 단어 목록)

#### Embedded Korean Word Lists
새로운 임베디드 한글 단어 파일들을 추가:

```go
// cmd/words/embedables/words/korean-common.json
{
  "metadata": {
    "name": "일반 한글 단어",
    "size": 1000,
    "packagedAt": "2024-01-01T00:00:00Z",
    "version": 1
  },
  "words": ["안녕", "세상", "컴퓨터", "프로그래밍", ...]
}

// cmd/words/embedables/words/korean-tech.json  
{
  "metadata": {
    "name": "한글 기술 용어",
    "size": 500,
    "packagedAt": "2024-01-01T00:00:00Z", 
    "version": 1
  },
  "words": ["데이터베이스", "알고리즘", "인터페이스", ...]
}

// cmd/words/embedables/sentences/korean-sentences.json
{
  "metadata": {
    "name": "한글 문장",
    "size": 100,
    "packagedAt": "2024-01-01T00:00:00Z",
    "version": 1
  },
  "words": ["안녕하세요. 반갑습니다.", "프로그래밍은 재미있습니다.", ...]
}
```

#### Korean Words Generator Extension
기존 `words.WordsGenerator`를 확장하여 한글 지원:

```go
// cmd/words/korean.go (새 파일)
package words

//go:embed embedables/words/korean-common.json
var koreanCommon string

//go:embed embedables/words/korean-tech.json  
var koreanTech string

//go:embed embedables/sentences/korean-sentences.json
var koreanSentences string

func addEmbeddedKoreanSources(sources map[string]WordSource) map[string]WordSource {
    // 한글 단어 소스들을 추가하는 함수
}

func NewKoreanGenerator(paths []string) WordsGenerator {
    // 한글 전용 생성기 생성
}
```

### 2. Menu System Extension (메뉴 시스템 확장)

#### Korean Test Settings
기존 테스트 설정 구조체들을 확장:

```go
// cmd/model.go에 추가
type KoreanTimerBasedTestSettings struct {
    TimerBasedTestSettings
    language string // "korean"
}

type KoreanWordCountBasedTestSettings struct {
    WordCountBasedTestSettings  
    language string // "korean"
}

type KoreanSentenceCountBasedTestSettings struct {
    SentenceCountBasedTestSettings
    language string // "korean"
}
```

#### Menu Structure Update
메인 메뉴에 한글 옵션 추가:

```go
// cmd/init.go 수정
func initMainMenu() MainMenu {
    config := ReadConfig()
    
    // 기존 영어 선택들
    timeBasedWordSelections := filterEnabledSelections(config)
    countBasedWordSelections := filterEnabledWordSelection(config)
    countBasedSentenceSelections := filterEnabledSentenceSelection(config)
    
    // 새로운 한글 선택들
    koreanTimeBasedSelections := filterEnabledKoreanSelections(config)
    koreanWordSelections := filterEnabledKoreanWordSelection(config)
    koreanSentenceSelections := filterEnabledKoreanSentenceSelection(config)
    
    return MainMenu{
        config: config,
        selections: []MainMenuSelection{
            // 영어 테스트들
            initTimerBasedTestSettings(config, timeBasedWordSelections),
            initWordCountBasedTestSettings(config, countBasedWordSelections),
            initSentenceCountBasedTestSettings(config, countBasedSentenceSelections),
            
            // 한글 테스트들  
            initKoreanTimerBasedTestSettings(config, koreanTimeBasedSelections),
            initKoreanWordCountBasedTestSettings(config, koreanWordSelections),
            initKoreanSentenceCountBasedTestSettings(config, koreanSentenceSelections),
            
            initConfigViewSelection(),
        },
        // 생성기들
        timeBasedGenerator:         words.NewGenerator(paths(timeBasedWordSelections)),
        wordCountGenerator:         words.NewGenerator(paths(countBasedWordSelections)),
        sentenceCountGenerator:     words.NewGenerator(paths(countBasedSentenceSelections)),
        koreanTimeBasedGenerator:   words.NewKoreanGenerator(paths(koreanTimeBasedSelections)),
        koreanWordCountGenerator:   words.NewKoreanGenerator(paths(koreanWordSelections)),
        koreanSentenceGenerator:    words.NewKoreanGenerator(paths(koreanSentenceSelections)),
    }
}
```

### 3. Korean Text Processing (한글 텍스트 처리)

#### Korean Character Handling
한글 문자 입력 및 처리를 위한 새로운 함수들:

```go
// cmd/korean.go (새 파일)
package cmd

import "unicode"

// 한글 문자인지 확인
func isKoreanChar(r rune) bool {
    return unicode.Is(unicode.Hangul, r)
}

// 한글 WPM 계산 (한글 특성 고려)
func calculateKoreanWPM(chars int, duration time.Duration) int {
    // 한글은 조합형 문자이므로 영어와 다른 계산 방식 적용
    // 일반적으로 한글 1글자 = 영어 2-3글자로 계산
    adjustedChars := float64(chars) * 2.5
    minutes := duration.Minutes()
    return int(adjustedChars / 5.0 / minutes) // 표준 WPM 공식 적용
}

// 한글 텍스트 정확도 계산
func calculateKoreanAccuracy(input, target []rune) float64 {
    // 한글 조합 문자의 특성을 고려한 정확도 계산
}
```

#### Korean Text Display
한글 텍스트 표시 및 줄바꿈 처리:

```go
// cmd/view.go 수정
func wrapKoreanStyledParagraph(paragraph string, lineLimit int) string {
    // 한글 문자의 폭을 고려한 줄바꿈 처리
    // 한글은 고정폭 문자이므로 영어와 다른 처리 필요
}

func (base *TestBase) koreanParagraphView(lineLimit int, styles Styles) string {
    // 한글 전용 문단 뷰 렌더링
}
```

### 4. Configuration Extension (설정 확장)

#### Korean Word Lists Configuration
설정 구조체에 한글 지원 추가:

```go
// cmd/model.go 수정
type Config struct {
    TestSettingCursors TestSettingCursors
    EmbededWordLists   []EmbededWordList
    WordLists          []WordList
    
    // 한글 지원 추가
    EmbeddedKoreanWordLists []EmbededWordList
    KoreanWordLists         []WordList
    
    LayoutFiles        []LayoutFile
    Layout             Layout
    Version            int
}

type EmbededWordList struct {
    Name        string
    IsSentences bool
    Enabled     bool
    Language    string // "english" 또는 "korean"
}
```

## Data Models

### Korean Word Source Structure
한글 단어 소스의 JSON 구조:

```json
{
  "metadata": {
    "name": "일반 한글 단어",
    "size": 1000,
    "packagedAt": "2024-01-01T00:00:00Z",
    "version": 1,
    "language": "korean",
    "category": "common"
  },
  "words": [
    "안녕",
    "세상", 
    "컴퓨터",
    "프로그래밍",
    "데이터베이스"
  ]
}
```

### Korean Test Results
한글 테스트 결과 구조:

```go
type KoreanResults struct {
    Results
    Language        string  // "korean"
    KoreanWPM      int     // 한글 특화 WPM
    CharacterCount int     // 한글 문자 수
    SyllableCount  int     // 음절 수
}
```

## Error Handling

### Korean Input Validation
한글 입력 검증 및 오류 처리:

1. **Invalid Korean Characters**: 지원하지 않는 한글 문자 입력 시 처리
2. **IME Input Handling**: 한글 입력기(IME) 상태 처리
3. **Incomplete Syllables**: 미완성 한글 음절 처리
4. **Mixed Language Input**: 한글/영어 혼합 입력 처리

### Error Recovery
```go
func handleKoreanInputError(input rune, expected rune) error {
    if !isKoreanChar(input) && isKoreanChar(expected) {
        return fmt.Errorf("expected Korean character, got: %c", input)
    }
    return nil
}
```

## Testing Strategy

### Unit Tests
1. **Korean Character Processing Tests**
   - 한글 문자 인식 테스트
   - 한글 WPM 계산 테스트
   - 한글 정확도 계산 테스트

2. **Korean Word Generation Tests**
   - 한글 단어 목록 로딩 테스트
   - 한글 문장 생성 테스트
   - 커스텀 한글 파일 처리 테스트

3. **Korean Menu Integration Tests**
   - 한글 메뉴 표시 테스트
   - 한글/영어 전환 테스트
   - 한글 설정 저장/로드 테스트

### Integration Tests
1. **End-to-End Korean Typing Tests**
   - 완전한 한글 타이핑 테스트 시나리오
   - 한글 결과 계산 및 표시 테스트
   - 한글 커스텀 파일 통합 테스트

### Test Data
```go
// cmd/korean_test.go
var koreanTestWords = []string{
    "안녕하세요",
    "프로그래밍", 
    "데이터베이스",
    "알고리즘",
}

var koreanTestSentences = []string{
    "안녕하세요. 반갑습니다.",
    "프로그래밍은 매우 재미있는 활동입니다.",
    "한글 타이핑 연습을 통해 실력을 향상시킬 수 있습니다.",
}
```

## Implementation Phases

### Phase 1: Core Korean Support
- 한글 단어 생성기 구현
- 기본 한글 단어 목록 추가
- 한글 문자 처리 함수 구현

### Phase 2: Menu Integration  
- 메인 메뉴에 한글 옵션 추가
- 한글 테스트 설정 구조체 구현
- 한글/영어 전환 기능 구현

### Phase 3: Advanced Features
- 한글 WPM 계산 최적화
- 커스텀 한글 파일 지원
- 한글 텍스트 표시 개선

### Phase 4: Testing & Polish
- 종합 테스트 수행
- 성능 최적화
- 사용자 경험 개선