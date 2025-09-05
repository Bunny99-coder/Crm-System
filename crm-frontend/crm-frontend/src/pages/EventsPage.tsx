// Replace the entire contents of src/pages/EventsPage.tsx with this.

import { useState, useEffect, useCallback } from 'react';
import { getEvents, deleteEvent } from '../services/apiService';
import type { Event } from '../services/apiService';
import { Table, Loader, Alert, Title, Button, Group, ActionIcon } from '@mantine/core';
import { IconPencil, IconTrash, IconAlertCircle } from '@tabler/icons-react';
import { useDisclosure } from '@mantine/hooks';
import CreateEventModal from '../components/CreateEventModal';

const EventsPage = () => {
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpened, { open: openModal, close: closeModal }] = useDisclosure(false);
  const [eventToEdit, setEventToEdit] = useState<Event | null>(null);

  const fetchEvents = useCallback(async () => {
    try {
      // Don't set loading to true on refetch, only on initial load
      const data = await getEvents();
      // --- THE FIX IS HERE ---
      // If data is null or undefined, default to an empty array.
      setEvents(data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load events. Are you logged in?');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchEvents();
  }, [fetchEvents]);

  const handleEditClick = (event: Event) => {
    setEventToEdit(event);
    openModal();
  };

  const handleCreateClick = () => {
    setEventToEdit(null);
    openModal();
  };
  
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure you want to cancel this event?')) {
      try {
        await deleteEvent(id);
        fetchEvents();
      } catch (err) {
        setError('Failed to delete event.');
      }
    }
  };

  if (loading) return <Loader color="blue" />;
  if (error) return <Alert color="red" title="Error" icon={<IconAlertCircle />}>{error}</Alert>;

  const rows = events.map((event) => (
    <Table.Tr key={event.id}>
      <Table.Td>{event.event_name}</Table.Td>
      <Table.Td>{new Date(event.start_time).toLocaleString()}</Table.Td>
      <Table.Td>{new Date(event.end_time).toLocaleString()}</Table.Td>
      <Table.Td>{event.location || 'N/A'}</Table.Td>
      <Table.Td>
        <Group gap="xs">
          <ActionIcon variant="subtle" onClick={() => handleEditClick(event)}><IconPencil size={16} /></ActionIcon>
          <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteClick(event.id)}><IconTrash size={16} /></ActionIcon>
        </Group>
      </Table.Td>
    </Table.Tr>
  ));

  return (
    <>
      <CreateEventModal
        opened={modalOpened}
        onClose={closeModal}
        onSuccess={fetchEvents}
        eventToEdit={eventToEdit}
      />

      <Group justify="space-between" mb="md">
        <Title order={2}>My Calendar / Events</Title>
        <Button onClick={handleCreateClick}>Schedule New Event</Button>
      </Group>

      <Table striped withTableBorder withColumnBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Event</Table.Th>
            <Table.Th>Start Time</Table.Th>
            <Table.Th>End Time</Table.Th>
            <Table.Th>Location</Table.Th>
            <Table.Th>Actions</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {rows.length > 0 ? (
            rows
          ) : (
            <Table.Tr>
              <Table.Td colSpan={5} align="center">No events found. Schedule one!</Table.Td>
            </Table.Tr>
          )}
        </Table.Tbody>
      </Table>
    </>
  );
};

export default EventsPage;