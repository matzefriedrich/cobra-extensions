### Unit Testing Guidelines

All unit tests in this project must follow these guidelines:

*   **AAA Principle**: Tests should be structured using the Arrange, Act, and Assert pattern.
    *   **Arrange**: Set up the conditions for the test (initialize objects, prepare inputs).
    *   **Act**: Execute the function or method being tested.
    *   **Assert**: Verify that the outcome is as expected.
*   **Naming Convention**: Test methods must follow the naming pattern: `Test_<Target>_<test_case_description_in_strict_snake_case>`.
    *   Example: `Test_commandDescriptor_bind_flags_correctly_maps_flags_to_cobra_command`
