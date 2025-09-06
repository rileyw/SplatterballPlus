# AGENTS.md

## Project Overview
SplatterballPlus is a multi-component game project with several submodules, including game logic, launcher, server, map editor, and helper libraries. It contains C#, C++, and resource files, organized into solution folders for modular development.

## Setup Commands
- Build all C# projects: Open `Magestorm.sln` in Visual Studio and build the solution.
- Build C++ components: Open `MageHook/MageHook.vcxproj` in Visual Studio and build.
- Restore NuGet packages: Visual Studio will prompt or use `nuget restore Magestorm.sln`.
- Install dependencies for C++: Ensure required libraries in `MageHook/Lib/` are present.

## Testing Instructions
- Run unit tests (if present) via Visual Studio Test Explorer for each C# project.
- Manual testing: Launch `MageLauncher` and `MageServer` executables after build.
- For C++: Use Visual Studio's debugger and run-time checks.

## Code Style Guidelines
- C#: Follow standard .NET conventions (PascalCase for types/methods, camelCase for variables).
- C++: Use consistent indentation and descriptive variable names.
- Keep code modular; use project folders for logical separation.

## PR Instructions
- Title format: `[SplatterballPlus] <Short Description>`
- Run all builds and tests before submitting PRs.
- Ensure no new warnings or errors in build output.
- Update or add tests for new/changed code.

## Security Considerations
- Do not commit secrets or sensitive data.
- Review dependencies for vulnerabilities before updating.

## Additional Notes
- Large datasets and resources are stored in `Content/` and subfolders.
- For subprojects, add tailored AGENTS.md files if agent-specific instructions are needed.
- Treat AGENTS.md as living documentation—update as project conventions evolve.
