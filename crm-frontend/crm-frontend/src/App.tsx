// Replace the contents of src/App.tsx
import { AppShell, Burger, Group, Text } from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { Routes, Route, Link } from 'react-router-dom';
import ContactsPage from './pages/ContactsPage';
import PropertiesPage from './pages/PropertiesPage'; // Import
import LeadsPage from './pages/LeadsPage';         // Import
import DealsPage from './pages/DealsPage';           // Import
import LoginPage from './pages/LoginPage';
import { useAuth } from './context/AuthContext';
import classes from './App.module.css';
import TasksPage from './pages/TasksPage';
import ViewContactPage from './pages/ViewContactPage';
import NotesPage from './pages/NotesPage';
import EventsPage from './pages/EventsPage';
import DashboardPage from './pages/DashboardPage';




const NavLink = ({ to, children }: { to: string; children: React.ReactNode }) => (
  <Link to={to} className={classes.link}>
    {children}
  </Link>
);

function App() {
  const [opened, { toggle }] = useDisclosure();
  const { isAuthenticated, logout } = useAuth();

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{ width: 250, breakpoint: 'sm', collapsed: { mobile: !opened } }}
      padding="md"
    >
      <AppShell.Header className={classes.header}>
        <Group h="100%" px="md">
          <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
          <Text size="xl" fw={700}>Real Estate CRM</Text>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        <NavLink to="/">Dashboard</NavLink>
        <NavLink to="/contacts">Contacts</NavLink>
        <NavLink to="/properties">Properties</NavLink>
        <NavLink to="/leads">Leads</NavLink>
        <NavLink to="/deals">Deals</NavLink>
        <NavLink to="/tasks">Tasks</NavLink> 
        <NavLink to="/notes">Notes</NavLink> 
        <NavLink to="/events">Events</NavLink> 

        <div style={{ marginTop: 'auto' }}>
          {isAuthenticated ? (
            <button className={classes.link} style={{ width: '100%', textAlign: 'left' }} onClick={logout}>Logout</button>
          ) : (
            <NavLink to="/login">Login</NavLink>
          )}
        </div>
      </AppShell.Navbar>

      <AppShell.Main>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/contacts" element={<ContactsPage />} />
          <Route path="/contacts/:contactId" element={<ViewContactPage />} /> 

          <Route path="/properties" element={<PropertiesPage />} />
          <Route path="/leads" element={<LeadsPage />} />
          <Route path="/deals" element={<DealsPage />} />
          <Route path="/tasks" element={<TasksPage />} /> 
          <Route path="/notes" element={<NotesPage />} /> 
          <Route path="/events" element={<EventsPage />} /> 


  <Route path="/" element={<DashboardPage />} /> 
        </Routes>
      </AppShell.Main>
    </AppShell>
  );
}

export default App;