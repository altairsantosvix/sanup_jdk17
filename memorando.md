MEMORANDO
📌 Para: Equipe Técnica e Gestão de Projetos
📌 De: [Seu Nome] – [Seu Cargo]
📌 Assunto: Estratégia para Migração de Serviços Java 11 para Java 17 e Adequação às Novas Regras


---

1. Contexto

Com a migração dos 140 serviços para Java 17, novas exigências foram introduzidas no pipeline de entrega, incluindo:

1️⃣ Análise de Segurança pelo Sysdig – Bloqueia o deploy caso sejam encontradas vulnerabilidades críticas ou altas em dependências.
2️⃣ Novas Regras do Sonar – Agora exige cobertura mínima de testes, eliminando a possibilidade de exclusão de código da análise (Sonar Exclusion).
3️⃣ Compatibilidade de Código Legado – Muitos serviços passaram para produção antes dessas regras e agora precisam ser ajustados.

Essas mudanças impactam diretamente os prazos de entrega e exigem uma estratégia para evitar bloqueios constantes na pipeline.


---

2. Principais Impactos

📌 2.1. Segurança – Bloqueios pelo Sysdig
🔹 Vulnerabilidades são atualizadas diariamente, o que pode fazer com que um serviço aprovado hoje seja bloqueado amanhã.
🔹 Muitas dependências precisam ser atualizadas, mas algumas são críticas e podem causar quebras inesperadas.

📌 2.2. Qualidade de Código – Novas Regras do Sonar
🔹 Serviços que usavam Sonar Exclusion agora precisam de testes unitários e de integração para aprovação na pipeline.
🔹 Código legado pode ser reprovado por baixa cobertura de testes, mesmo sem alterações recentes.

📌 2.3. Código Legado – Risco de Incompatibilidade
🔹 Dependências antigas podem não ser compatíveis com Java 17 e precisarão de ajustes.
🔹 Atualizações podem impactar funcionalidades que não possuem cobertura de testes, aumentando o risco de falhas em produção.


---

3. Estratégia para Mitigação

Para garantir continuidade das entregas e evitar bloqueios inesperados, propomos as seguintes ações:

📌 3.1. Segurança – Redução de Impacto do Sysdig
✅ Criar um pipeline de análise antecipada de vulnerabilidades, permitindo correções antes da etapa final de entrega.
✅ Priorizar dependências LTS (Long-Term Support) para reduzir a necessidade de atualizações frequentes.
✅ Definir critérios para exceções justificadas, permitindo que vulnerabilidades não exploráveis na aplicação não bloqueiem a pipeline.
✅ Implementar uma rotina de atualização contínua, evitando acúmulo de problemas.

📌 3.2. Testes – Adequação ao Novo Sonar
✅ Focar primeiro nos serviços mais críticos, garantindo que aplicações essenciais atendam aos novos requisitos rapidamente.
✅ Criar testes de integração básicos para cobrir os principais fluxos sem grande impacto no desenvolvimento.
✅ Rodar os testes em uma pipeline separada, permitindo ajustes antes do deploy final.

📌 3.3. Código Legado – Ajustes Mínimos e Estratégicos
✅ Mapear os serviços mais problemáticos e priorizar sua migração gradual.
✅ Evitar grandes refatorações e focar apenas nos ajustes necessários para compatibilidade com Java 17.
✅ Sempre que possível, utilizar ferramentas de compatibilidade (Shims, Backwards Compatibility Layers) para facilitar a transição.


---

4. Conclusão

A nova pipeline traz benefícios em segurança e qualidade, mas também riscos operacionais. Para garantir entregas sem impacto significativo nos prazos, será essencial:

🔹 Implementar um pipeline de análise antecipada para segurança.
🔹 Priorizar testes automáticos mínimos para atender às novas regras do Sonar.
🔹 Focar em ajustes mínimos para compatibilidade de código legado, sem grandes refatorações.

A colaboração de todos será fundamental para garantir uma migração eficiente e segura.

📌 Ações imediatas:
✅ Criar um grupo de revisão para mapear os serviços mais críticos.
✅ Implementar uma rotina de atualização semanal/mensal para dependências.
✅ Ajustar a pipeline para permitir testes e validações antecipadas.

Para dúvidas ou sugestões, estou à disposição.

Atenciosamente,
[Seu Nome]
[Seu Cargo]

