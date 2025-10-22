---
description: Generate innovative app ideas for Next.js + Firebase applications and save to files
---

Generate 10 innovative app ideas that can be built with Next.js, TypeScript, Tailwind, and Firebase.

Requirements for each idea:
- Solves a real problem for a specific audience
- Completable in 2-5 days with Claude
- Can be built with Next.js App Router + Firebase (Firestore + Auth)
- Has clear value proposition
- Not oversaturated in the market
- Deployable as a Docker container to GCP Cloud Run

For each idea provide:

1. **App Name**: Creative, memorable name
2. **Tagline**: One sentence description
3. **Problem**: What specific problem does it solve?
4. **Target Audience**: Who would use this? (be specific)
5. **Core MVP Features** (3-5):
   - Feature 1
   - Feature 2
   - Feature 3
   - (etc.)
6. **Technical Complexity**: Simple / Medium / Complex
7. **Time Estimate**: 2-5 days
8. **Unique Value**: Why would someone use this over alternatives?
9. **Monetization Potential**: Free tier + paid features, subscription, etc.
10. **Example Use Cases**: 2-3 concrete examples

Focus on practical, useful applications that real people would want to use.

Consider these categories:
- Productivity tools
- Personal finance
- Health & wellness
- Education & learning
- Creative tools
- Social & community
- Business utilities
- Developer tools

Make the ideas specific and actionable.

## Output Format

After generating all 10 ideas, save EACH idea to a separate markdown file in `/Users/sloan/code/app-ideas/docs/app-ideas/` with the filename format: `{app-name-kebab-case}.md`

For example:
- `habit-tracker-pro.md`
- `recipe-cost-calculator.md`
- `meeting-cost-analyzer.md`

Each file should contain the complete idea in markdown format with all sections above.

At the end, provide a summary list of all generated ideas with their filenames.
