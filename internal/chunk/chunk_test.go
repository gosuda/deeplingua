package chunk_test

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"gosuda.org/deeplingua/internal/chunk"
)

func TestChunkMarkdown(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name: "Original test case",
			input: `# Test Markdown

This is a paragraph with some text.

## Code Block Test

` + "```go" + `
func example() {
	fmt.Println("Hello, World!")
}
` + "```" + `

Another paragraph after the code block.

- List item 1
- List item 2

> This is a blockquote.

Final paragraph with a [link](https://example.com).`,
		},
		{
			name: "Long Original test case",
			input: strings.Repeat(`# Test Markdown

This is a paragraph with some text.

## Code Block Test

`+"```go"+`
func example() {
	fmt.Println("Hello, World!")
}
`+"```"+`

Another paragraph after the code block.

- List item 1
- List item 2

> This is a blockquote.

Final paragraph with a [link](https://example.com).`, 4000),
		},
		{
			name: "Long paragraphs and multiple code blocks",
			input: `# Long Content Test

This is a very long paragraph that should exceed the token limit. It contains a lot of text to ensure that the chunkMarkdown function correctly handles long content. We want to make sure that it splits the content appropriately without breaking the markdown structure. This paragraph goes on and on with more text to really push the limits of our function.

## First Code Block

` + "```python" + `
def long_function():
    print("This is a long function")
    for i in range(100):
        print(f"Iteration {i}")
    return "Done"
` + "```" + `

Another paragraph here to separate the code blocks. We'll add more text to make it longer and ensure proper chunking behavior.

## Second Code Block

` + "```javascript" + `
function anotherLongFunction() {
    console.log("This is another long function");
    for (let i = 0; i < 100; i++) {
        console.log('Iteration ${i}');
    }
    return "Done";
}
` + "```" + `

Final paragraph after multiple code blocks to test handling of mixed content types.`,
		},
		{
			name: "Nested lists and block quotes",
			input: `# Nested Structures Test

1. First level item
   - Nested unordered item
   - Another nested item
     1. Deep nested ordered item
     2. Another deep nested item
2. Second level item
   > This is a block quote inside a list
   > It continues for multiple lines
   > To test proper handling of nested structures

- Unordered list item
  > Another block quote
  > With multiple lines

  And some text after the quote

3. Third level item
   ` + "```" + `
   This is a code block inside a list item
   It should be handled correctly
   ` + "```" + `

Final paragraph to conclude the nested structures test.`,
		},
		{
			name: "Headers, horizontal rules, and tables",
			input: `# Main Header

## Subheader 1

Some content under subheader 1.

### Sub-subheader

More content here.

---

## Subheader 2

Content under subheader 2.

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Row 1, Col 1 | Row 1, Col 2 | Row 1, Col 3 |
| Row 2, Col 1 | Row 2, Col 2 | Row 2, Col 3 |
| Row 3, Col 1 | Row 3, Col 2 | Row 3, Col 3 |

#### Deep nested header

* List item 1
* List item 2

---

# Another Main Header

Final paragraph to conclude the headers and tables test.`,
		},
		{
			name: "Very long content 1",
			input: strings.Repeat("This is a long sentence that will be repeated many times to create a very large content. ", 500) +
				"\n\n" + strings.Repeat("Another long sentence with different words to increase variety in the content. ", 500),
		},
		{
			name: "Very long content 2",
			input: strings.Repeat("\n\n가가가가가가가가가가가This is a long sentence that will be repeated many times to create a very large content.\n\n ", 500) +
				"\n\n" + strings.Repeat("Another long sentence with different words to increase variety in the content. ", 500),
		},
		{
			name:  "Very long content 3",
			input: strings.Repeat("\n\na", 10000),
		},
		{
			name:  "Very long content 4",
			input: strings.Repeat("가", 10000),
		},
		{
			name:  "Very long content 5",
			input: strings.Repeat("<think>\nOkay, let's tackle this problem step by step. It's a bit intimidating with all those logarithms and exponents, but I'll take it one piece at a time. The problem states:\n\nLet \\( a > 1 \\) and \\( x > 1 \\) satisfy:\n\\[\n\\log_a(\\log_a(\\log_a 2) + \\log_a 24 - 128) = 128\n\\]\nand\n\\[\n\\log_a(\\log_a x) = 256\n\\]\nFind the remainder when \\( x \\) is divided by 1000.\n\nFirst, let's parse the problem. We have two equations involving logarithms with base \\( a \\), and we need to find \\( x \\) modulo 1000. To do this, we probably need to find the value of \\( a \\) first, using the first equation, and then use that to solve for \\( x \\) in the second equation.\n\nStarting with the first equation:\n\\[\n\\log_a(\\log_a(\\log_a 2) + \\log_a 24 - 128) = 128\n\\]\nLet me rewrite this step by step. Remember that \\( \\log_b c = d \\) is equivalent to \\( b^d = c \\). So, applying that here, if I let the inner expression be \\( y \\), then:\n\\[\n\\log_a(y) = 128 \\implies a^{128} = y\n\\]\nWhere \\( y = \\log_a(\\log_a 2) + \\log_a 24 - 128 \\). So:\n\\[\na^{128} = \\log_a(\\log_a 2) + \\log_a 24 - 128\n\\]\nHmm, this looks complicated. Let's see if we can simplify the terms involving logarithms. Let's first compute \\( \\log_a(\\log_a 2) \\) and \\( \\log_a 24 \\).\n\nLet me denote \\( \\log_a 2 \\) as some variable to make it easier. Let's let \\( b = \\log_a 2 \\). Then, by definition, \\( a^b = 2 \\).\n\nSimilarly, \\( \\log_a 24 = \\log_a (24) \\), which I can leave as is for now. So substituting back into the equation:\n\\[\na^{128} = \\log_a(b) + \\log_a 24 - 128\n\\]\nBut \\( b = \\log_a 2 \\), so \\( \\log_a(b) = \\log_a(\\log_a 2) \\). Wait, that's the same as the original term. Maybe I need a different approach.\n\nAlternatively, let's work step by step, converting each logarithm into exponent form. Let's start from the outside and work our way in.\n\nGiven:\n\\[\n\\log_a(\\log_a(\\log_a 2) + \\log_a 24 - 128) = 128\n\\]\nAs before, converting the outer logarithm:\n\\[\na^{128} = \\log_a(\\log_a 2) + \\log_a 24 - 128\n\\]\nLet me denote \\( C = \\log_a(\\log_a 2) + \\log_a 24 \\). Then the equation becomes:\n\\[\na^{128} = C - 128 \\implies C = a^{128} + 128\n\\]\nSo:\n\\[\n\\log_a(\\log_a 2) + \\log_a 24 = a^{128} + 128\n\\]\nThis still seems complex. Maybe we can express all terms in terms of \\( \\log_a 2 \\) and \\( \\log_a 3 \\), since 24 is 2^3 * 3.\n\nYes, 24 factors into \\( 2^3 \\times 3 \\). Therefore:\n\\[\n\\log_a 24 = \\log_a (2^3 \\times 3) = 3\\log_a 2 + \\log_a 3\n\\]\nSo now, substituting back into the equation:\n\\[\n\\log_a(\\log_a 2) + 3\\log_a 2 + \\log_a 3 = a^{128} + 128\n\\]\nHmm, this might not immediately help. Let me see if there's a different approach. Perhaps assuming that \\( a \\) is a power of 2? Or some integer? The problem doesn't specify that \\( a \\) is an integer, but given the exponents are large (128, 256), maybe \\( a \\) is 2? Let's test that.\n\nAssume \\( a = 2 \\). Let's check if this satisfies the first equation.\n\nFirst, compute \\( \\log_2(\\log_2(\\log_2 2) + \\log_2 24 - 128) \\).\n\nStarting from the innermost log:\n\n\\( \\log_2 2 = 1 \\).\n\nThen, \\( \\log_2(\\log_2 2) = \\log_2 1 = 0 \\).\n\nNext, \\( \\log_2 24 \\). Since 24 is 2^3 * 3, \\( \\log_2 24 = \\log_2 (2^3 \\times 3) = 3 + \\log_2 3 \\approx 3 + 1.58496 = 4.58496 \\).\n\nSo substituting back into the expression inside the logarithm:\n\n\\( 0 + 4.58496 - 128 = -123.41504 \\).\n\nThen, \\( \\log_2(-123.41504) \\) is undefined because logarithm of a negative number is not real. So \\( a = 2 \\) doesn't work.\n\nOkay, so \\( a \\) is not 2. Maybe \\( a \\) is 4? Let's try \\( a = 4 \\).\n\nCompute \\( \\log_4(\\log_4(\\log_4 2) + \\log_4 24 - 128) \\).\n\nFirst, innermost log: \\( \\log_4 2 \\). Since \\( 4^{1/2} = 2 \\), so \\( \\log_4 2 = 1/2 \\).\n\nThen, \\( \\log_4(\\log_4 2) = \\log_4(1/2) = \\log_4(2^{-1}) = -1 \\times \\log_4 2 = -1/2 \\).\n\nNext, \\( \\log_4 24 \\). 24 is 4 * 6 = 4^1 * 6, so \\( \\log_4 24 = 1 + \\log_4 6 \\). Then \\( \\log_4 6 = \\log_4 (4^{1} \\times 1.5) = 1 + \\log_4 1.5 \\approx 1 + 0.2925 = 1.2925 \\). Therefore, \\( \\log_4 24 \\approx 1 + 1.2925 = 2.2925 \\).\n\nPutting it all into the expression inside the logarithm:\n\n\\( -0.5 + 2.2925 - 128 = -126.2075 \\).\n\nAgain, trying to take the log base 4 of a negative number, which is undefined. So \\( a = 4 \\) is not working either.\n\nHmm. Maybe \\( a \\) is a larger number? Let's think differently.\n\nLet me recall that the outermost logarithm equals 128. So the argument of that logarithm must be \\( a^{128} \\). Let me write that again:\n\nIf \\( \\log_a(Z) = 128 \\), then \\( Z = a^{128} \\). So here, \\( Z = \\log_a(\\log_a 2) + \\log_a 24 - 128 \\).\n\nTherefore:\n\\[\n\\log_a(\\log_a 2) + \\log_a 24 - 128 = a^{128}\n\\]\nOr rearranged:\n\\[\n\\log_a(\\log_a 2) + \\log_a 24 = a^{128} + 128\n\\]\nThis equation seems really hard to solve directly. Maybe we need to make a substitution. Let me denote \\( t = \\log_a 2 \\). Then \\( a^t = 2 \\). Also, \\( \\log_a 24 = \\log_a (2^3 \\times 3) = 3\\log_a 2 + \\log_a 3 = 3t + \\log_a 3 \\). Let me also denote \\( s = \\log_a 3 \\). Then \\( a^s = 3 \\).\n\nSo substituting into the equation:\n\\[\n\\log_a(t) + 3t + s = a^{128} + 128\n\\]\nBut \\( t = \\log_a 2 \\), so \\( \\log_a(t) = \\log_a(\\log_a 2) \\). Hmm, this seems circular. Maybe we can express \\( a \\) in terms of \\( t \\) and \\( s \\)?\n\nGiven that \\( a^t = 2 \\) and \\( a^s = 3 \\), so \\( a = 2^{1/t} = 3^{1/s} \\). Therefore, \\( 2^{1/t} = 3^{1/s} \\). Taking logarithms base 2:\n\\[\n\\frac{1}{t} = \\frac{\\log_2 3}{s} \\implies s = t \\log_2 3\n\\]\nSo \\( s = t \\times \\log_2 3 \\approx t \\times 1.58496 \\).\n\nTherefore, we can write \\( s \\) in terms of \\( t \\). So substituting back into the equation:\n\\[\n\\log_a(t) + 3t + t \\log_2 3 = a^{128} + 128\n\\]\nBut \\( a = 2^{1/t} \\), so \\( \\log_a(t) = \\frac{\\log_2 t}{\\log_2 a} = \\frac{\\log_2 t}{1/t} = t \\log_2 t \\).\n\nTherefore, substituting back:\n\\[\nt \\log_2 t + 3t + t \\log_2 3 = a^{128} + 128\n\\]\nSimplify the left-hand side:\n\\[\nt \\log_2 t + t(3 + \\log_2 3) = a^{128} + 128\n\\]\nFactor out \\( t \\):\n\\[\nt \\left( \\log_2 t + 3 + \\log_2 3 \\right) = a^{128} + 128\n\\]\nBut \\( a = 2^{1/t} \\), so \\( a^{128} = (2^{1/t})^{128} = 2^{128/t} \\). Therefore:\n\\[\nt \\left( \\log_2 t + 3 + \\log_2 3 \\right) = 2^{128/t} + 128\n\\]\nThis equation is still quite complicated. Let's see if we can find a value of \\( t \\) that makes this equation hold. Since \\( a > 1 \\), and \\( t = \\log_a 2 \\), if \\( a > 2 \\), then \\( t < 1 \\). If \\( 1 < a < 2 \\), then \\( t > 1 \\). But given that the exponent 128 is involved, maybe \\( a \\) is a number such that \\( 2^{128/t} \\) is manageable?\n\nAlternatively, suppose that \\( 128/t \\) is an integer. Let's assume \\( 128/t = k \\), where \\( k \\) is a positive integer. Then \\( t = 128/k \\), and \\( a = 2^{1/t} = 2^{k/128} \\). Let's see if this assumption leads us anywhere.\n\nSubstituting \\( t = 128/k \\) into the left-hand side:\n\\[\n\\frac{128}{k} \\left( \\log_2 \\left( \\frac{128}{k} \\right) + 3 + \\log_2 3 \\right)\n\\]\nAnd the right-hand side becomes:\n\\[\n2^{k} + 128\n\\]\nSo we have:\n\\[\n\\frac{128}{k} \\left( \\log_2 \\left( \\frac{128}{k} \\right) + 3 + \\log_2 3 \\right) = 2^{k} + 128\n\\]\nHmm, perhaps trying small integer values for \\( k \\). Let's try \\( k = 7 \\). Then \\( t = 128/7 ≈ 18.2857 \\), which seems large, but let's check.\n\nLeft-hand side:\n\\[\n\\frac{128}{7} \\left( \\log_2 \\left( \\frac{128}{7} \\right) + 3 + \\log_2 3 \\right)\n\\]\nCalculate \\( \\frac{128}{7} ≈ 18.2857 \\)\n\\( \\log_2(128/7) = \\log_2(128) - \\log_2(7) = 7 - approximately 2.8074 ≈ 4.1926 \\)\nThen, \\( 4.1926 + 3 + 1.58496 ≈ 8.77756 \\)\nMultiply by \\( 128/7 ≈ 18.2857 \\times 8.77756 ≈ 160.6 \\)\nRight-hand side: \\( 2^7 + 128 = 128 + 128 = 256 \\). So 160.6 ≠ 256. Not equal.\n\nTry \\( k = 8 \\). Then \\( t = 128/8 = 16 \\)\nLeft-hand side:\n\\( 128/8 = 16 \\)\nInside the brackets: \\( \\log_2(16) + 3 + log_2 3 = 4 + 3 + 1.58496 ≈ 8.58496 \\)\nMultiply by 16: 16 * 8.58496 ≈ 137.359\nRHS: \\( 2^8 + 128 = 256 + 128 = 384 \\). Not equal.\n\nk=6: t=128/6≈21.333\nLHS: 128/6 ≈21.333\nInside: log2(128/6)=log2(21.333)≈4.415 +3 +1.58496≈9.0\nMultiply by 21.333≈21.333*9≈192\nRHS: 2^6 +128=64+128=192. Hey, that's equal!\n\nWait, hold on. Let me check:\n\nIf k=6, then t=128/6≈21.3333\n\nCompute LHS:\n\nFirst, compute log2(128/6). Since 128/6 ≈21.3333, log2(21.3333). Since 2^4=16, 2^5=32, so log2(21.3333)=4 + log2(1.3333)=4 + 0.415≈4.415.\n\nThen, 4.415 + 3 + log2 3 ≈4.415 +3 +1.58496≈9.0.\n\nMultiply by 128/6≈21.3333: 21.3333 *9≈192.\n\nRHS: 2^6 +128=64 +128=192. Perfect, it matches!\n\nSo, k=6 is a solution. Therefore, t=128/6=64/3≈21.3333, so a=2^{k/128}=2^{6/128}=2^{3/64}.\n\nWait, hold on. Wait, if k=6, then t=128/k≈21.3333, but a=2^{k/128}=2^{6/128}=2^{3/64}. Let me confirm:\n\nSince we had earlier that a=2^{1/t}, and t=128/k, so 1/t=k/128, so a=2^{k/128}. Thus, with k=6, a=2^{6/128}=2^{3/64}.\n\nSo, a=2^{3/64}. Let me compute that value if necessary, but maybe we can keep it as 2^{3/64} for now.\n\nTherefore, we found that k=6 gives a solution. So, a=2^{3/64}. Let me check if this is the only solution. Since the equation might have multiple solutions, but given the context of the problem and the answer expecting an integer modulo 1000, likely there's a unique solution here.\n\nTherefore, a=2^{3/64}.\n\nNow, moving on to the second equation:\n\\[\n\\log_a(\\log_a x) = 256\n\\]\nAgain, let's convert this logarithmic equation to its exponential form. If \\( \\log_a(Y) = 256 \\), then Y = a^{256}. Here, Y = \\( \\log_a x \\), so:\n\\[\n\\log_a x = a^{256}\n\\]\nThen, converting again:\n\\[\nx = a^{a^{256}}\n\\]\nSo, x is a tower of exponents: a raised to the power of a^256. Given that a=2^{3/64}, this becomes:\n\\[\nx = \\left(2^{3/64}\\right)^{(2^{3/64})^{256}}\n\\]\nSimplify the exponents step by step.\n\nFirst, compute the exponent in the power of a: (2^{3/64})^{256}\n\nRecall that (2^{b})^{c} = 2^{b*c}. So:\n\n(2^{3/64})^{256} = 2^{(3/64)*256} = 2^{(3*256)/64} = 2^{(768)/64} = 2^{12} = 4096.\n\nTherefore, the exponent simplifies to 4096. Therefore:\n\nx = (2^{3/64})^{4096} = 2^{(3/64)*4096} = 2^{(3*4096)/64}\n\nCalculate 4096 divided by 64: 4096 /64 = 64. Because 64*64=4096.\n\nTherefore, (3*4096)/64 = 3*64 = 192.\n\nTherefore, x=2^{192}.\n\nSo, x=2^{192}. The problem asks for the remainder when x is divided by 1000. So, we need to compute 2^{192} mod 1000.\n\nCalculating 2^192 mod 1000. That's a huge exponent, but we can use Euler's theorem or observe patterns in the powers of 2 modulo 1000.\n\nFirst, note that 1000=8*125, and since 2 and 125 are coprime (gcd(2,125)=1), we can use Euler's theorem for modulo 125. Euler's totient function φ(125)=100, since 125 is 5^3, so φ(5^3)=5^3 -5^2=125-25=100.\n\nTherefore, 2^100 ≡1 mod 125.\n\nSimilarly, modulo 8: 2^3=8≡0 mod8, so any higher power of 2 will be 0 mod8. So 2^192 mod8=0.\n\nTherefore, using the Chinese Remainder Theorem, we can compute 2^192 mod1000 by finding a number congruent to 0 mod8 and congruent to 2^192 mod125. Then combine them.\n\nFirst, compute 2^192 mod125.\n\nSince 2^100≡1 mod125, as per Euler's theorem. So 192=100*1 +92, so 2^192= (2^100)^1 *2^92≡1*2^92 mod125.\n\nSo need to compute 2^92 mod125.\n\nCompute 2^92 mod125.\n\nWe can compute this using exponentiation by squaring.\n\nFirst, compute 2^1=2 mod125\n\n2^2=4\n\n2^4=(2^2)^2=16\n\n2^8=(2^4)^2=256≡6 mod125\n\n2^16=(2^8)^2=6^2=36 mod125\n\n2^32=(2^16)^2=36^2=1296≡1296-10*125=1296-1250=46 mod125\n\n2^64=(2^32)^2=46^2=2116≡2116-16*125=2116-2000=116 mod125\n\n2^64≡116 mod125.\n\nNow, 2^92=2^64 *2^16 *2^8 *2^4. Wait, 64+16+8+4=92. Let's check:\n\nWait 64 + 16 + 8 + 4=92? 64+16=80, +8=88, +4=92. Yes.\n\nTherefore, 2^92≡(116)*(36)*(6)*(16) mod125.\n\nCompute step by step:\n\nFirst, 116*36 mod125:\n\n116*36=4176. 4176 divided by125: 125*33=4125, 4176-4125=51. So 116*36≡51 mod125.\n\nNext, 51*6=306. 306-2*125=306-250=56. So 51*6≡56 mod125.\n\nThen, 56*16=896. 896-7*125=896-875=21. So 56*16≡21 mod125.\n\nTherefore, 2^92≡21 mod125.\n\nThus, 2^192≡21 mod125.\n\nTherefore, 2^192≡21 mod125 and ≡0 mod8. We need to find a number x such that:\n\nx≡0 mod8\n\nx≡21 mod125\n\nFind x ≡21 mod125 and x≡0 mod8.\n\nLet x=125k +21. We need 125k +21≡0 mod8.\n\n125 mod8=5 (since 8*15=120, 125-120=5). Therefore:\n\n5k +21≡0 mod8\n\n21 mod8=5, so:\n\n5k +5≡0 mod8 ⇒5k≡-5≡3 mod8.\n\nMultiply both sides by inverse of 5 mod8. Since 5*5=25≡1 mod8, inverse of 5 is 5.\n\nThus, k≡3*5=15≡7 mod8.\n\nTherefore, k=8m +7 for integer m. Therefore:\n\nx=125*(8m +7)+21=1000m +875 +21=1000m +896.\n\nTherefore, the smallest positive solution is x≡896 mod1000.\n\nTherefore, 2^192≡896 mod1000.\n\nTherefore, the remainder when x is divided by1000 is 896.\n\nBut wait, let me double-check my computations to be sure.\n\nFirst, compute 2^192 mod125=21, as above. Then:\n\nFind x≡21 mod125 and x≡0 mod8.\n\nSet x=125k +21. Then 125k +21≡5k +5≡0 mod8. So 5k≡-5≡3 mod8.\n\nInverse of 5 mod8 is 5, since 5*5=25≡1 mod8. So k≡3*5=15≡7 mod8. Hence, k=8m +7. Then x=125*(8m +7)+21=1000m +875 +21=1000m +896. So indeed x≡896 mod1000.\n\nTherefore, the remainder is 896.\n\nHence, the answer is 896.\n</think>\n\nFinal Answer: The final answer is $\\boxed{896}", 10),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("Running test case: %s\n", tc.name)
			chunks := chunk.ChunkMarkdown(tc.input)
			fmt.Printf("Number of chunks: %d\n", len(chunks))

			for i := range chunks {
				if !utf8.ValidString(chunks[i]) {
					t.Errorf("Invalid UTF-8 string in chunk %d for test case '%s'", i, tc.name)
				}
			}

			joined := strings.Join(chunks, "")

			if joined != tc.input {
				t.Errorf("Joined chunks do not match original input for test case '%s'.\nExpected length: %d\nGot length: %d", tc.name, len(tc.input), len(joined))

				// Find the first difference
				for i := 0; i < len(tc.input) && i < len(joined); i++ {
					if tc.input[i] != joined[i] {
						fmt.Printf("First difference at index %d. Expected: %q, Got: %q\n", i, tc.input[i], joined[i])
						break
					}
				}
			} else {
				fmt.Printf("Test case '%s' passed successfully\n", tc.name)
			}

			// Additional check for very long content
			if tc.name == "Very long content" {
				if len(chunks) < 2 {
					t.Errorf("Expected multiple chunks for very long content, got %d chunks", len(chunks))
				} else {
					fmt.Printf("Very long content split into %d chunks as expected\n", len(chunks))
				}
			}

			fmt.Println() // Add a blank line between test cases
		})
	}
}
