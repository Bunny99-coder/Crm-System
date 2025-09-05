// Replace the entire contents of src/pages/LeadsPage.tsx with this.
import { useState, useEffect, useCallback } from 'react';
import { getLeads, deleteLead } from '../services/apiService';
import type { Lead } from '../services/apiService';
import { Table, Loader, Alert, Title, Button, Group, ActionIcon } from '@mantine/core';
import { IconPencil, IconTrash, IconAlertCircle } from '@tabler/icons-react';
import { useDisclosure } from '@mantine/hooks';
import CreateLeadModal from '../components/CreateLeadModal';
import { useAuth } from '../context/AuthContext';
import { useData } from '../context/DataContext';

const LeadsPage = () => {
  const { userRole } = useAuth();
  const isReception = userRole === 2; // Role ID for Reception

  const [leads, setLeads] = useState<Lead[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpened, { open: openModal, close: closeModal }] = useDisclosure(false);
  const [leadToEdit, setLeadToEdit] = useState<Lead | null>(null);

  const { userMap, contactMap, loading: dataLoading } = useData();

  const fetchLeads = useCallback(async () => {
    try {
      setLoading(true);
      const data = await getLeads();
      setLeads(data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load leads. Are you logged in?');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchLeads();
  }, [fetchLeads]);

  const handleEditClick = (lead: Lead) => {
    setLeadToEdit(lead);
    openModal();
  };
  const handleCreateClick = () => {
    setLeadToEdit(null);
    openModal();
  };
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this lead?')) {
      try {
        await deleteLead(id);
        fetchLeads();
      } catch (err) {
        setError('Failed to delete lead.');
      }
    }
  };

  if ((loading || dataLoading) && leads.length === 0) return <Loader color="blue" />;
  if (error) return <Alert color="red" title="Error" icon={<IconAlertCircle />}>{error}</Alert>;

  const rows = leads.map((lead) => (
    <Table.Tr key={lead.id}>
      <Table.Td>{lead.id}</Table.Td>
      <Table.Td>{contactMap.get(lead.contact_id) || `ID: ${lead.contact_id}`}</Table.Td>
      <Table.Td>{userMap.get(lead.assigned_to) || `ID: ${lead.assigned_to}`}</Table.Td>
      <Table.Td>{lead.status_id}</Table.Td>
      <Table.Td>{new Date(lead.created_at).toLocaleDateString()}</Table.Td>
      {/* Conditionally render the Actions column only for Reception */}
      {isReception && (
        <Table.Td>
          <Group gap="xs">
            <ActionIcon variant="subtle" onClick={() => handleEditClick(lead)}><IconPencil size={16} /></ActionIcon>
            <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteClick(lead.id)}><IconTrash size={16} /></ActionIcon>
          </Group>
        </Table.Td>
      )}
    </Table.Tr>
  ));

  return (
    <>
      <CreateLeadModal opened={modalOpened} onClose={closeModal} onSuccess={fetchLeads} leadToEdit={leadToEdit} />
      <Group justify="space-between" mb="md">
        <Title order={2}>Leads</Title>
        {/* Conditionally render the "Create" button only for Reception */}
        {isReception && <Button onClick={handleCreateClick}>Create New Lead</Button>}
      </Group>
      <Table striped withTableBorder withColumnBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Lead ID</Table.Th>
            <Table.Th>Contact</Table.Th>
            <Table.Th>Assigned To</Table.Th>
            <Table.Th>Status ID</Table.Th>
            <Table.Th>Date Created</Table.Th>
            {/* Conditionally render the Actions header */}
            {isReception && <Table.Th>Actions</Table.Th>}
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {leads.length > 0 ? rows : (
             <Table.Tr><Table.Td colSpan={isReception ? 6 : 5} align="center">No leads found.</Table.Td></Table.Tr>
          )}
        </Table.Tbody>
      </Table>
    </>
  );
};

export default LeadsPage;