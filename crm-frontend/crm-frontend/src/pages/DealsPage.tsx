// Replace the entire contents of src/pages/DealsPage.tsx with this.

import { useState, useEffect, useCallback } from 'react';
import { getDeals, deleteDeal } from '../services/apiService';
import type { Deal } from '../services/apiService';
import { Table, Loader, Alert, Title, Button, Group, ActionIcon, Badge } from '@mantine/core';
import { IconPencil, IconTrash, IconAlertCircle } from '@tabler/icons-react';
import { useDisclosure } from '@mantine/hooks';
import CreateDealModal from '../components/CreateDealModal';
import { useAuth } from '../context/AuthContext';
import { useData } from '../context/DataContext';

const DealsPage = () => {
  const { userRole } = useAuth();
  const isSalesAgent = userRole === 1; // Correctly identify if the user is a Sales Agent

  const [deals, setDeals] = useState<Deal[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpened, { open: openModal, close: closeModal }] = useDisclosure(false);
  const [dealToEdit, setDealToEdit] = useState<Deal | null>(null);

  const { propertyMap, leadMap, loading: dataLoading } = useData();

  const fetchDeals = useCallback(async () => {
    try {
      setLoading(true);
      const data = await getDeals();
      setDeals(data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load deals. Please try again.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchDeals();
  }, [fetchDeals]);

  const handleEditClick = (deal: Deal) => {
    setDealToEdit(deal);
    openModal();
  };
  const handleCreateClick = () => {
    setDealToEdit(null);
    openModal();
  };
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this deal?')) {
      try {
        await deleteDeal(id);
        fetchDeals();
      } catch (err) {
        setError('Failed to delete deal.');
      }
    }
  };

  if ((loading || dataLoading) && deals.length === 0) return <Loader color="blue" />;
  if (error) return <Alert color="red" title="Error" icon={<IconAlertCircle />}>{error}</Alert>;

  const rows = deals.map((deal) => (
    <Table.Tr key={deal.id}>
      <Table.Td>{deal.id}</Table.Td>
      <Table.Td>{leadMap.get(deal.lead_id) || `Lead ID: ${deal.lead_id}`}</Table.Td>
      <Table.Td>{propertyMap.get(deal.property_id) || `Property ID: ${deal.property_id}`}</Table.Td>
      <Table.Td>
        <Badge color={deal.deal_status === 'Closed-Won' ? 'green' : 'gray'}>
            {deal.deal_status}
        </Badge>
      </Table.Td>
      <Table.Td>${deal.deal_amount.toLocaleString()}</Table.Td>
      {/* Conditionally render the Actions column ONLY for Sales Agents */}
      {isSalesAgent && (
        <Table.Td>
          <Group gap="xs">
            <ActionIcon variant="subtle" onClick={() => handleEditClick(deal)}><IconPencil size={16} /></ActionIcon>
            <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteClick(deal.id)}><IconTrash size={16} /></ActionIcon>
          </Group>
        </Table.Td>
      )}
    </Table.Tr>
  ));

  return (
    <>
      <CreateDealModal opened={modalOpened} onClose={closeModal} onSuccess={fetchDeals} dealToEdit={dealToEdit} />
      <Group justify="space-between" mb="md">
        <Title order={2}>Deals</Title>
        {/* Conditionally render the "Create" button ONLY for Sales Agents */}
        {isSalesAgent && <Button onClick={handleCreateClick}>Create New Deal</Button>}
      </Group>
      <Table striped withTableBorder withColumnBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Deal ID</Table.Th>
            <Table.Th>Lead</Table.Th>
            <Table.Th>Property</Table.Th>
            <Table.Th>Status</Table.Th>
            <Table.Th>Amount</Table.Th>
            {/* Conditionally render the Actions header ONLY for Sales Agents */}
            {isSalesAgent && <Table.Th>Actions</Table.Th>}
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {deals.length > 0 ? rows : (
            <Table.Tr><Table.Td colSpan={isSalesAgent ? 6 : 5} align="center">No deals found.</Table.Td></Table.Tr>
          )}
        </Table.Tbody>
      </Table>
    </>
  );
};

export default DealsPage;