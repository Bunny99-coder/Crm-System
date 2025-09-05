// Replace the entire contents of src/pages/PropertiesPage.tsx
import { useState, useEffect, useCallback } from 'react';
import { getProperties, deleteProperty } from '../services/apiService';
import type { Property } from '../services/apiService';
import { Table, Loader, Alert, Title, Button, Group, ActionIcon } from '@mantine/core';
import { IconPencil, IconTrash, IconAlertCircle } from '@tabler/icons-react';
import { useDisclosure } from '@mantine/hooks';
import CreatePropertyModal from '../components/CreatePropertyModal';
import { useAuth } from '../context/AuthContext';

const PropertiesPage = () => {
  const { userRole } = useAuth();
  const isReception = userRole === 2;

  const [properties, setProperties] = useState<Property[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpened, { open: openModal, close: closeModal }] = useDisclosure(false);
  const [propertyToEdit, setPropertyToEdit] = useState<Property | null>(null);

  const fetchProperties = useCallback(async () => {
    try {
      setLoading(true);
      const data = await getProperties();
      setProperties(data || []);
      setError(null);
    } catch (err) { setError('Failed to load properties. Are you logged in?'); }
    finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchProperties(); }, [fetchProperties]);

  const handleEditClick = (property: Property) => { setPropertyToEdit(property); openModal(); };
  const handleCreateClick = () => { setPropertyToEdit(null); openModal(); };
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure?')) {
      try {
        await deleteProperty(id);
        fetchProperties();
      } catch (err) { setError('Failed to delete property.'); }
    }
  };

  if (loading && properties.length === 0) return <Loader color="blue" />;
  if (error) return <Alert color="red" title="Error" icon={<IconAlertCircle />}>{error}</Alert>;

  const rows = properties.map((property) => (
    <Table.Tr key={property.id}>
      <Table.Td>{property.name}</Table.Td>
      <Table.Td>{property.unit_no || 'N/A'}</Table.Td>
      <Table.Td>{property.status}</Table.Td>
      <Table.Td>${property.price.toLocaleString()}</Table.Td>
      {isReception && (
        <Table.Td>
          <Group gap="xs">
            <ActionIcon variant="subtle" onClick={() => handleEditClick(property)}><IconPencil size={16} /></ActionIcon>
            <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteClick(property.id)}><IconTrash size={16} /></ActionIcon>
          </Group>
        </Table.Td>
      )}
    </Table.Tr>
  ));

  return (
    <>
      <CreatePropertyModal opened={modalOpened} onClose={closeModal} onSuccess={fetchProperties} propertyToEdit={propertyToEdit} />
      <Group justify="space-between" mb="md">
        <Title order={2}>Properties</Title>
        {isReception && <Button onClick={handleCreateClick}>Create New Property</Button>}
      </Group>
      <Table striped withTableBorder withColumnBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Property Name</Table.Th>
            <Table.Th>Unit No.</Table.Th>
            <Table.Th>Status</Table.Th>
            <Table.Th>Price</Table.Th>
            {isReception && <Table.Th>Actions</Table.Th>}
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
           {properties.length > 0 ? rows : (<Table.Tr><Table.Td colSpan={isReception ? 5: 4} align="center">No properties found.</Table.Td></Table.Tr>)}
        </Table.Tbody>
      </Table>
    </>
  );
};
export default PropertiesPage;