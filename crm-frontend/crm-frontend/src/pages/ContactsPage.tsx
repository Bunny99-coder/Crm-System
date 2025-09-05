// Replace the entire contents of src/pages/ContactsPage.tsx
import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { getContacts, deleteContact } from '../services/apiService';
import type { Contact } from '../services/apiService';
import { Table, Loader, Alert, Title, Button, Group, ActionIcon } from '@mantine/core';
import { IconPencil, IconTrash, IconAlertCircle } from '@tabler/icons-react';
import { useDisclosure } from '@mantine/hooks';
import CreateContactModal from '../components/CreateContactModal';
import { useAuth } from '../context/AuthContext';

const ContactsPage = () => {
  const { userRole } = useAuth();
  const isReception = userRole === 2;

  const [contacts, setContacts] = useState<Contact[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpened, { open: openModal, close: closeModal }] = useDisclosure(false);
  const [contactToEdit, setContactToEdit] = useState<Contact | null>(null);

  const fetchContacts = useCallback(async () => {
    try {
      setLoading(true);
      const data = await getContacts();
      setContacts(data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load contacts. Are you logged in?');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchContacts();
  }, [fetchContacts]);

  const handleEditClick = (contact: Contact) => {
    setContactToEdit(contact);
    openModal();
  };
  const handleCreateClick = () => {
    setContactToEdit(null);
    openModal();
  };
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure?')) {
      try {
        await deleteContact(id);
        fetchContacts();
      } catch (err) { setError("Failed to delete contact."); }
    }
  };

  if (loading && contacts.length === 0) return <Loader color="blue" />;
  if (error) return <Alert color="red" title="Error" icon={<IconAlertCircle />}>{error}</Alert>;

  const rows = contacts.map((contact) => (
    <Table.Tr key={contact.id}>
      <Table.Td>
        <Link to={`/contacts/${contact.id}`} style={{ textDecoration: 'none' }}>
          {contact.first_name} {contact.last_name}
        </Link>
      </Table.Td>
      <Table.Td>{contact.email || 'N/A'}</Table.Td>
      <Table.Td>{contact.primary_phone}</Table.Td>
      {isReception && (
        <Table.Td>
          <Group gap="xs">
            <ActionIcon variant="subtle" onClick={() => handleEditClick(contact)}><IconPencil size={16} /></ActionIcon>
            <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteClick(contact.id)}><IconTrash size={16} /></ActionIcon>
          </Group>
        </Table.Td>
      )}
    </Table.Tr>
  ));

  return (
    <>
      <CreateContactModal opened={modalOpened} onClose={closeModal} onSuccess={fetchContacts} contactToEdit={contactToEdit} />
      <Group justify="space-between" mb="md">
        <Title order={2}>Contacts</Title>
        {isReception && <Button onClick={handleCreateClick}>Create New Contact</Button>}
      </Group>
      <Table striped withTableBorder withColumnBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Email</Table.Th>
            <Table.Th>Phone</Table.Th>
            {isReception && <Table.Th>Actions</Table.Th>}
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {contacts.length > 0 ? rows : (
            <Table.Tr><Table.Td colSpan={isReception ? 4 : 3} align="center">No contacts found.</Table.Td></Table.Tr>
          )}
        </Table.Tbody>
      </Table>
    </>
  );
};
export default ContactsPage;