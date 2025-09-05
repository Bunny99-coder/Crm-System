// Replace the entire contents of src/pages/ViewContactPage.tsx

import { useParams } from 'react-router-dom';
import { useState, useEffect, useCallback } from 'react';
import { getContactById, getCommLogsForContact, deleteCommLog } from '../services/apiService';
import type { Contact, CommLog } from '../services/apiService';
import { Title, Paper, Text, Loader, Alert, Table, Group, Button, ActionIcon } from '@mantine/core';
import { IconAlertCircle, IconPencil, IconTrash } from '@tabler/icons-react';
import { useDisclosure } from '@mantine/hooks';
import CreateCommLogModal from '../components/CreateCommLogModal';

const ViewContactPage = () => {
  const { contactId } = useParams<{ contactId: string }>();
  const [contact, setContact] = useState<Contact | null>(null);
  // Initialize 'logs' as an empty array, which is safe for .map()
  const [logs, setLogs] = useState<CommLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpened, { open: openModal, close: closeModal }] = useDisclosure(false);
  const [logToEdit, setLogToEdit] = useState<CommLog | null>(null);
  
  const contactIdNum = Number(contactId);

  const fetchData = useCallback(async () => {
    if (!contactIdNum) {
      setError("Invalid Contact ID.");
      setLoading(false);
      return;
    }

    try {
      // Don't set loading to true here, it's already true from the initial state
      const [contactData, logsData] = await Promise.all([
        getContactById(contactIdNum),
        getCommLogsForContact(contactIdNum),
      ]);
      setContact(contactData);
      setLogs(logsData || []); // <-- Safety check: if logsData is null/undefined, use an empty array
      setError(null);
    } catch (err) {
      setError('Failed to load contact details.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [contactIdNum]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);


  // --- HANDLERS ---
  const handleEditClick = (log: CommLog) => {
    setLogToEdit(log);
    openModal();
  };
  const handleCreateClick = () => {
    setLogToEdit(null);
    openModal();
  };
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this log entry?')) {
      try {
        await deleteCommLog(id);
        fetchData(); // Refetch all data
      } catch (err) {
        setError('Failed to delete log entry.');
      }
    }
  };


  // --- RENDER LOGIC ---
  if (loading) return <Loader />;
  if (error || !contact) return <Alert color="red" title="Error" icon={<IconAlertCircle />}>{error || 'Contact not found.'}</Alert>;

  // By the time we get here, 'logs' is guaranteed to be an array (even if empty).
  const logRows = logs.map((log) => (
    <Table.Tr key={log.id}>
      <Table.Td>{new Date(log.interaction_date).toLocaleString()}</Table.Td>
      <Table.Td>{log.interaction_type}</Table.Td>
      <Table.Td>{log.notes || 'N/A'}</Table.Td>
      <Table.Td>
        <Group gap="xs">
          <ActionIcon variant="subtle" onClick={() => handleEditClick(log)}><IconPencil size={16} /></ActionIcon>
          <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteClick(log.id)}><IconTrash size={16} /></ActionIcon>
        </Group>
      </Table.Td>
    </Table.Tr>
  ));

  return (
    <>
      <CreateCommLogModal
        opened={modalOpened}
        onClose={closeModal}
        onSuccess={fetchData}
        logToEdit={logToEdit}
        contactId={contactIdNum}
      />
      <Paper shadow="xs" p="xl" mb="xl">
        <Title order={2}>{contact.first_name} {contact.last_name}</Title>
        <Text>Email: {contact.email || 'N/A'}</Text>
        <Text>Phone: {contact.primary_phone}</Text>
      </Paper>
      <Group justify="space-between" mb="md">
        <Title order={3}>Communication History</Title>
        <Button onClick={handleCreateClick}>Add Log Entry</Button>
      </Group>
      <Table>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Date</Table.Th>
            <Table.Th>Type</Table.Th>
            <Table.Th>Notes</Table.Th>
            <Table.Th>Actions</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>{logRows.length > 0 ? logRows : <Table.Tr><Table.Td colSpan={4} align="center">No log entries found.</Table.Td></Table.Tr>}</Table.Tbody>
      </Table>
    </>
  );
};

export default ViewContactPage;