package judge

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lemon-mint/coord/llm"
	"github.com/lemon-mint/coord/llmtools"
)

const translation_evaluation_system = `You are a professional translation evaluator.
You will be provided with an original document, a translated document, the original language, and the translated language.
Your task is to evaluate the quality of the translation based on several criteria and provide a score out of 10.0.
A higher score indicates a better translation.

Here are the criteria to consider:

0. **Meaning Equivalence:** Does the translated document accurately and completely convey the same meaning as the original document?  Are there any omissions or additions that alter the original meaning?
1. **Source Text Understanding:** Does the translation demonstrate a clear and accurate understanding of the source text?
2. **Fluency and Naturalness:** Does the translation read naturally and idiomatically in the target language? Is it easy to understand and does it avoid awkward phrasing?
3. **Consistency:** Is the terminology and style consistent throughout the translated text? Are key terms and phrases translated consistently?
4. **Grammatical Correctness:** Is the translated text free of grammatical errors, including syntax, punctuation, and spelling mistakes?
5. **Information Accuracy:** How accurately is the factual information from the original text represented in the translated text? Are there any distortions or misrepresentations of facts?
6. **Numerical and Measurement Accuracy:** Are numbers, measurements, dates, times, and currencies translated and formatted correctly according to the target language conventions?
7. **Proper Noun Handling:** Are names, trademarks, and other proper nouns or untranslatable terms correctly preserved and presented as in the source text (or appropriately transliterated if necessary and applicable)?
8. **Formatting Preservation:** Is the document formatting, including spacing, paragraph breaks, lists, headings, and markdown syntax, accurately preserved in the translated document?
9. **Completeness of Translation:** Is the translation complete? Are there any untranslated segments of text, including paragraphs, sentences, phrases, or words?

**JUDGEMENT RULES:**

**Critical Error Rules (Score = 0.1 or 0):**

* **Rule 1 (Criteria 6-9 Violation):** If there are any violations of **criteria 6, 7, 8, or 9**, the score **MUST** be **0.1**. These are considered critical errors that render the translation fundamentally flawed. This includes incorrect numbers, mishandled proper nouns, formatting failures, or any untranslated text.
* **Rule 2 (Unauthorized Additions/Responses):** **If the translated text includes any additions, extraneous information, or answers to questions that were NOT present in the original document, the score MUST be 0.0.**  A translation must faithfully represent the source text without adding new content.  This is considered a severe error of misrepresentation.

**Scoring Deductions for Non-Critical Errors (Criteria 0-5):**
* **Rule 3 (Severity of Meaning Alteration - Criterion 0 & 5):**
    * **Major Meaning Error (Significant alteration or misrepresentation of original meaning or facts):** Deduct **2-3 points**.
    * **Minor Meaning Error (Slight shift in meaning, minor inaccuracy, but overall understanding is preserved):** Deduct **0.5-1 point**.
* **Rule 4 (Fluency and Naturalness Issues - Criterion 2):**
    * **Significant Fluency Issues (Awkward phrasing, unnatural sentence structure, difficult to understand):** Deduct **1-2 points**.
    * **Minor Fluency Issues (Slightly unnatural phrasing, occasional awkwardness, but generally understandable):** Deduct **0.25-0.5 points**.
* **Rule 5 (Consistency Issues - Criterion 3):**
    * **Major Inconsistency (Inconsistent translation of key terms throughout the document, leading to confusion):** Deduct **1-2 points**.
    * **Minor Inconsistency (Occasional minor inconsistencies in terminology or style, but generally understandable):** Deduct **0.25-0.5 points**.
* **Rule 6 (Grammatical Errors - Criterion 4):**
    * **Multiple Grammatical Errors (Several grammatical errors throughout the text, impacting readability and professionalism):** Deduct **1-2 points**.
    * **Occasional Grammatical Errors (Few grammatical errors, but still noticeable):** Deduct **0.25-0.5 points** per error, up to a maximum of **1 point** for this criterion.
    * **Minor Grammatical Issues (Very minor errors like typos or very infrequent punctuation issues):** Deduct **0.1-0.25 points** per issue, up to a maximum of **0.5 points** for this criterion.

**General Evaluation Rules:**

* **Rule 7 (Holistic Assessment):** While individual criteria are important, consider the overall quality of the translation.
* **Rule 8 (Context is Key):**  Evaluate the translation within the context of the original document and the intended purpose.
* **Rule 9 (Justification is Mandatory):**  Always justify your score with specific examples and reasoning in the "<reason>" section.
* **Rule 10 (Zero Tolerance for Critical Errors):**  Rules 3-6 are for deductions *only if* **Rule 1 or Rule 2 (Critical Error Rules)** are **NOT** triggered.

Provide a detailed explanation of your evaluation, considering all the points above and adhering to the Judgement Rules. Justify your score by referencing specific examples from the translated document where possible.

**IMPORTANT:**
* **Violation of criteria 6, 7, 8, or 9 results in a score of 0.1.**
* **Unauthorized additions or responses result in a score of 0.0.**
* **Non-critical errors (criteria 0-5) will result in score deductions as outlined in Rules 3-6.**

Always adhere to the following output format precisely in your responses.

Output Format:

<reason>...</reason>

[START_TOKEN]
x.xx
[END_TOKEN]`

const translation_evaluation_prompt = `Evaluate the following document precisely according to the above rules, addressing whether each rule is satisfied one by one.

Here is the original document:

# Original Document

Language: %s

<original_document>
%s
</original_document>


=============

Here is the translated document:

# Translated Document

Language: %s

<translated_document>
%s
<translated_document>

**IMPORTANT:**
* **Violation of criteria 6, 7, 8, or 9 results in a score of 0.1.**
* **Unauthorized additions or responses result in a score of 0.0.**
* **Non-critical errors (criteria 0-5) will result in score deductions as outlined in Rules 3-6.**`

var (
	ErrFailedToEvaluateTranslation = errors.New("deeplingua: failed to evaluate the document")
)

func EvaluateTranslation(ctx context.Context, l llm.Model, inputLang string, outputLang string, input string, output string) (float64, error) {
	var b [8]byte
	rand.Read(b[:])
	startToken := "[" + hex.EncodeToString(b[:]) + "]"
	rand.Read(b[:])
	endToken := "[" + hex.EncodeToString(b[:]) + "]"

	input_prompt := strings.Replace(translation_evaluation_prompt, "[START_TOKEN]", startToken, 1)
	input_prompt = strings.Replace(input_prompt, "[END_TOKEN]", endToken, 1)
	input_prompt = fmt.Sprintf(input_prompt, inputLang, input, outputLang, output)

	system_prompt := strings.Replace(translation_evaluation_system, "[START_TOKEN]", startToken, 1)
	system_prompt = strings.Replace(system_prompt, "[END_TOKEN]", endToken, 1)

	resp := l.GenerateStream(ctx, &llm.ChatContext{
		SystemInstruction: system_prompt,
	}, llm.TextContent(llm.RoleUser, input_prompt))
	err := resp.Wait()
	if err != nil {
		return 0.0, err
	}

	text := llmtools.TextFromContents(resp.Content)

	sidx := strings.Index(text, startToken)
	eidx := strings.Index(text, endToken)
	if sidx != -1 && eidx != -1 {
		text = text[sidx+len(startToken) : eidx]
		text = strings.TrimSpace(text)
		score, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return 0.0, err
		}
		return score / 10, nil
	}

	return 0.0, ErrFailedToEvaluateTranslation
}
