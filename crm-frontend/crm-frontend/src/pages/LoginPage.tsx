// Replace the contents of src/pages/LoginPage.tsx
import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { TextInput, PasswordInput, Button, Paper, Title, Container, Alert } from '@mantine/core';
import { IconAlertCircle } from '@tabler/icons-react';

const LoginPage = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const { login } = useAuth();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await login({ username, password });
      // Navigation is now handled by the AuthContext
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container size={420} my={40}>
      <Title ta="center">Welcome Back!</Title>
      
      <Paper withBorder shadow="md" p={30} mt={30} radius="md">
        <form onSubmit={handleSubmit}>
          <TextInput
            label="Username"
            placeholder="Your username"
            value={username}
            onChange={(event) => setUsername(event.currentTarget.value)}
            required
          />
          <PasswordInput
            label="Password"
            placeholder="Your password"
            value={password}
            onChange={(event) => setPassword(event.currentTarget.value)}
            required
            mt="md"
          />
          <Button fullWidth mt="xl" type="submit" loading={loading}>
            Sign in
          </Button>
          {error && (
            <Alert
              variant="light"
              color="red"
              title="Login Error"
              icon={<IconAlertCircle />}
              mt="md"
            >
              {error}
            </Alert>
          )}
        </form>
      </Paper>
    </Container>
  );
};

export default LoginPage;