package openai

import "fmt"

func buildSystemPrompt() string {
	return `Você é um analista sênior de software escrevendo descrições para cards JIRA.

Regras:
- Escreva em Português (BR)
- Use linguagem não-técnica, acessível para QA e PO
- Foque no impacto funcional, não em detalhes de implementação
- Máximo ~200 palavras
- Formato obrigatório:

**Resumo:** [1-2 frases descrevendo a mudança principal]

**Mudanças Realizadas:**
- [bullet point descrevendo mudança funcional]
- [bullet point descrevendo mudança funcional]

**Impacto Funcional:** [Como isso afeta o usuário final ou o fluxo do sistema]`
}

func buildFewShotExamples() []Message {
	return []Message{
		{
			Role: "user",
			Content: `Descreva as mudanças deste commit:

diff --git a/src/auth/login.ts b/src/auth/login.ts
--- a/src/auth/login.ts
+++ b/src/auth/login.ts
@@ -15,6 +15,12 @@ export async function login(email: string, password: string) {
+  if (loginAttempts >= 5) {
+    await lockAccount(email, 30);
+    throw new AccountLockedError('Conta bloqueada por 30 minutos');
+  }
+
+  loginAttempts++;`,
		},
		{
			Role: "assistant",
			Content: `**Resumo:** Adicionado mecanismo de bloqueio temporário de conta após tentativas excessivas de login.

**Mudanças Realizadas:**
- Após 5 tentativas de login sem sucesso, a conta é bloqueada automaticamente por 30 minutos
- O sistema agora contabiliza cada tentativa de login realizada

**Impacto Funcional:** Usuários que errarem a senha 5 vezes terão a conta temporariamente bloqueada, aumentando a segurança contra tentativas de acesso não autorizadas.`,
		},
		{
			Role: "user",
			Content: `Gere uma descrição para card JIRA baseado neste PR:

diff --git a/src/checkout/payment.ts b/src/checkout/payment.ts
--- a/src/checkout/payment.ts
+++ b/src/checkout/payment.ts
@@ -8,3 +8,15 @@ const paymentMethods = ['credit_card', 'debit_card'];
+const paymentMethods = ['credit_card', 'debit_card', 'pix'];
+
+async function processPixPayment(order: Order) {
+  const pixCode = await generatePixCode(order.total);
+  return { method: 'pix', code: pixCode, expiresIn: '30m' };
+}`,
		},
		{
			Role: "assistant",
			Content: `**Resumo:** Adicionado PIX como nova opção de pagamento no checkout.

**Mudanças Realizadas:**
- PIX foi incluído como método de pagamento disponível no checkout
- Ao escolher PIX, um código é gerado automaticamente com validade de 30 minutos

**Impacto Funcional:** Clientes agora podem pagar via PIX durante o checkout, oferecendo uma opção de pagamento instantâneo além de cartão de crédito e débito.`,
		},
	}
}

func buildUserPrompt(diff, analysisType string, commitMessages []string, prTitle, prDescription string) string {
	var prefix string
	switch analysisType {
	case "single_commit":
		prefix = "Descreva as mudanças deste commit:"
	case "commit_range":
		prefix = "Descreva as mudanças consolidadas destes commits:"
	case "pull_request":
		prefix = "Gere uma descrição para card JIRA baseado neste PR:"
	default:
		prefix = "Descreva as mudanças:"
	}

	var context string
	if prTitle != "" {
		context += fmt.Sprintf("\nTítulo do PR: %s", prTitle)
	}
	if prDescription != "" {
		context += fmt.Sprintf("\nDescrição do PR: %s", prDescription)
	}
	if len(commitMessages) > 0 {
		context += "\nMensagens de commit:"
		for _, msg := range commitMessages {
			context += fmt.Sprintf("\n- %s", msg)
		}
	}

	return fmt.Sprintf("%s%s\n\n%s", prefix, context, diff)
}

func buildRefinePrompt(original, instruction string) string {
	return fmt.Sprintf(`Refine a seguinte descrição de card JIRA conforme a instrução.

Descrição original:
%s

Instrução de refinamento: %s

Mantenha o mesmo formato (Resumo, Mudanças Realizadas, Impacto Funcional) e escreva em Português (BR).`, original, instruction)
}

type Message struct {
	Role    string
	Content string
}
