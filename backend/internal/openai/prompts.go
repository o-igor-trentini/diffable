package openai

import "fmt"

func buildSystemPrompt() string {
	return buildSystemPromptForLevel("functional")
}

func buildSystemPromptForLevel(level string) string {
	switch level {
	case "technical":
		return `Você é um engenheiro de software sênior escrevendo descrições técnicas detalhadas.

Regras:
- Escreva em Português (BR)
- Use linguagem técnica, referenciando código, padrões de projeto e trade-offs
- Inclua detalhes de implementação relevantes (arquivos, funções, classes)
- Máximo ~300 palavras
- Formato obrigatório:

**Resumo Técnico:** [1-2 frases descrevendo a mudança principal com contexto técnico]

**Mudanças Implementadas:**
- [bullet point com detalhes técnicos: arquivo, função, padrão usado]
- [bullet point com detalhes técnicos]

**Decisões Técnicas:** [Trade-offs, padrões escolhidos, motivação técnica]

**Impacto:** [Efeitos em performance, segurança, manutenibilidade ou arquitetura]`

	case "executive":
		return `Você é um analista de negócios escrevendo resumos executivos concisos.

Regras:
- Escreva em Português (BR)
- Use linguagem de negócios, sem termos técnicos
- Foque no valor entregue e impacto para o negócio
- Máximo 2-3 frases
- Formato obrigatório:

**Resumo Executivo:** [2-3 frases descrevendo o que mudou, por que importa e qual o impacto para o negócio/usuário final]`

	case "qa_detailed":
		return `Você é um analista de qualidade sênior escrevendo descrições detalhadas para validação de QA em cards JIRA.

Regras:
- Escreva em Português (BR)
- Analise o diff de forma completa e minuciosa
- Priorize completude sobre brevidade — não há limite rígido de palavras
- Identifique mudanças de banco de dados, APIs, regras de negócio, fluxos afetados e cenários de teste
- Omita seções que não se aplicam ao diff analisado
- Formato obrigatório (use apenas as seções aplicáveis):

**Contexto:** [O que a mudança faz e por que foi necessária]

**Mudanças no Banco de Dados:** [Tabelas, colunas, migrações, índices, constraints]

**Mudanças de API/Integrações:** [Endpoints novos/alterados, serviços externos, novos campos de request/response]

**Regras de Negócio:** [Validações, lógica condicional, restrições, limites]

**Fluxos Afetados:**
- Caminho feliz: [fluxo principal esperado]
- Alternativo: [fluxos secundários]
- Caso de borda: [situações limítrofes ou excepcionais]

**Cenários de Teste Sugeridos:**
- [Cenário com dados de entrada e resultado esperado]
- [Cenário com dados de entrada e resultado esperado]

**Observações:** [Retrocompatibilidade, dependências, riscos, pontos de atenção]`

	default: // "functional"
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
}

func buildFewShotExamplesForLevel(level string) []Message {
	switch level {
	case "qa_detailed":
		return []Message{
			{
				Role: "user",
				Content: `Gere uma descrição detalhada para QA validar este PR:

diff --git a/src/integration/cnh-nexus.ts b/src/integration/cnh-nexus.ts
--- a/src/integration/cnh-nexus.ts
+++ b/src/integration/cnh-nexus.ts
@@ -1,8 +1,25 @@
+import { NexusClient } from '@nexus/sdk';
+
+interface CNHData {
+  numero: string;
+  dataEmissao: string | null;
+  categoria: string;
+}
+
+async function fetchCNHFromNexus(cpf: string): Promise<CNHData> {
+  const client = new NexusClient();
+  const response = await client.consultarCNH(cpf);
+  return {
+    numero: response.documentNumber,
+    dataEmissao: response.issueDate ?? null,
+    categoria: response.category,
+  };
+}
diff --git a/db/migrations/20240115_add_cnh_fields.sql b/db/migrations/20240115_add_cnh_fields.sql
+ALTER TABLE clientes ADD COLUMN cnh_data_emissao DATE NULL;
+ALTER TABLE clientes ADD COLUMN cnh_categoria VARCHAR(5) NULL;`,
			},
			{
				Role: "assistant",
				Content: `**Contexto:** Integração com a API Nexus para consulta de dados da CNH do cliente. O campo de data de emissão é opcional e pode vir vazio dependendo do estado emissor.

**Mudanças no Banco de Dados:**
- Tabela ` + "`clientes`" + `: adicionadas colunas ` + "`cnh_data_emissao`" + ` (DATE, nullable) e ` + "`cnh_categoria`" + ` (VARCHAR(5), nullable)
- Migração: ` + "`20240115_add_cnh_fields.sql`" + `

**Mudanças de API/Integrações:**
- Nova integração com API Nexus via SDK (` + "`@nexus/sdk`" + `)
- Endpoint consumido: ` + "`consultarCNH(cpf)`" + ` — retorna número, data de emissão e categoria
- Campo ` + "`dataEmissao`" + ` pode retornar ` + "`null`" + ` (tratado com nullish coalescing)

**Regras de Negócio:**
- Data de emissão da CNH é opcional (nullable) — alguns estados não fornecem essa informação
- Categoria da CNH é limitada a 5 caracteres (ex: "AB", "AE")

**Fluxos Afetados:**
- Caminho feliz: CPF válido → API Nexus retorna dados completos → dados salvos na tabela clientes
- Alternativo: CPF válido, mas data de emissão ausente → campo salvo como NULL
- Caso de borda: CPF inválido ou API Nexus indisponível → tratamento de erro necessário

**Cenários de Teste Sugeridos:**
- Enviar CPF válido com todos os campos preenchidos → verificar persistência correta
- Enviar CPF válido sem data de emissão → verificar que ` + "`cnh_data_emissao`" + ` é NULL
- Enviar CPF inválido → verificar tratamento de erro
- Simular API Nexus indisponível → verificar fallback/erro adequado

**Observações:**
- As colunas são nullable, sem impacto em registros existentes
- Dependência externa: SDK Nexus deve estar configurado no ambiente de QA`,
			},
		}
	default:
		return buildFewShotExamples()
	}
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

func buildUserPrompt(diff, analysisType string, commitMessages []string, prTitle, prDescription, userContext string) string {
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
	if userContext != "" {
		context += fmt.Sprintf("\n\nContexto adicional fornecido pelo desenvolvedor:\n%s", userContext)
	}
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
