# AI Agents for Development

GoTemplate includes three specialized AI agents designed to streamline your development workflow by breaking down complex tasks into manageable steps while following hexagonal architecture principles.

## 🔍 Clarifier Agent

![Clarifier](./images/clarifier.png "GoTemplate Clarifier")

**Purpose**: Understands and defines problems clearly before any solution is discussed.

**Usage**: Reference `@.github/agents/go-template_clarifier.agent.md` when you need to:
- Clarify vague requirements
- Define problem scope and boundaries
- Identify which GoTemplate components are involved
- Establish measurable success criteria

## 📋 Task Designer Agent

![Task Designer](./images/task-designer.png "GoTemplate Task Designer")

**Purpose**: Breaks down clarified problems into executable, reviewable tasks following GoTemplate's architecture.

**Usage**: Reference `@.github/agents/go-template_task_designer.agent.md` when you need to:
- Create task breakdowns from problem definitions
- Plan development work with proper dependencies
- Design tasks that respect hexagonal architecture boundaries
- Reference appropriate Makefile targets for code generation

## ⚡ Executor Agent

![Executor](./images/executor.png "GoTemplate Executor")

**Purpose**: Executes tasks precisely as designed without changing scope or adding features.

**Usage**: Reference `@.github/agents/go-template_executor.agent.md` when you need to:
- Implement tasks exactly as specified
- Follow Definition of Done criteria strictly
- Execute one task at a time with proper completion validation
- Maintain focus on requirements without scope creep

## Quick Start with Agents

1. **Start with Clarifier**: `@go-template_clarifier.agent.md "I want to add user authentication"`
2. **Design Tasks**: `@go-template_task_designer.agent.md` (use clarifier output)
3. **Execute Tasks**: `@go-template_executor.agent.md` (use task designer output)

This three-agent workflow ensures clear requirements, proper planning, and precise execution while maintaining GoTemplate's architectural integrity.