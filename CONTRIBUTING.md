# Contributing

Thank you for helping improve Marionette.

## Development flow

1. Open an issue or draft PR to discuss intent.
2. Implement in small, reviewable commits.
3. Update docs when behavior, APIs, or UX changes.
4. Run tests before opening PR.

## Documentation responsibilities

When changing user-visible behavior, update at least one of:

- `docs/site/getting-started.md`
- `docs/site/framework-overview.md`
- `docs/site/cookbook/`
- `docs/site/components/`

When changing architecture or conventions, update:

- `docs/architecture/*.md`
- `docs/ui-component-guidelines.md`

## Pull request checklist

- [ ] Behavior change is explained in PR description.
- [ ] Documentation is updated or marked N/A.
- [ ] Tests are added/updated where appropriate.
- [ ] No new `.ts` / `.tsx` files were added without approval.

## Ownership and review

- Ownership rules are defined in `.github/CODEOWNERS`.
- Use the PR template in `.github/pull_request_template.md`.
- Keep user docs (`docs/site/`) and engineering docs (`docs/`) updates in separate commits when possible.

## Documentation authoring

- Use `docs/site/templates/page-template.md` for new user-facing pages.
- Link related pages at the bottom to improve navigation.
- Prefer short, task-oriented pages over long mixed-reference pages.
