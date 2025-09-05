// Replace the entire contents of src/pages/NotesPage.tsx with this.
import { useState, useEffect, useCallback } from 'react';
import { getNotesForUser, deleteNote } from '../services/apiService';
import type { Note } from '../services/apiService';
import { Title, Button, Group, ActionIcon, Card, Text, SimpleGrid, Loader, Alert } from '@mantine/core';
import { IconPencil, IconTrash, IconAlertCircle } from '@tabler/icons-react';
import { useDisclosure } from '@mantine/hooks';
import CreateNoteModal from '../components/CreateNoteModal';
import { useAuth } from '../context/AuthContext';
import { getUserIdFromToken } from '../util/auth';

const NotesPage = () => {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpened, { open: openModal, close: closeModal }] = useDisclosure(false);
  const [noteToEdit, setNoteToEdit] = useState<Note | null>(null);
  
  const { token } = useAuth();

  const fetchNotes = useCallback(async () => {
    const userId = getUserIdFromToken(token);
    if (!userId) {
      setError("You must be logged in to view notes.");
      setLoading(false);
      return;
    }
    
    try {
      setLoading(true);
      const data = await getNotesForUser(userId);
      setNotes(data || []);
      setError(null);
    } catch (err) {
      setError("Could not fetch notes.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    fetchNotes();
  }, [fetchNotes]);

  const handleEditClick = (note: Note) => {
    setNoteToEdit(note);
    openModal();
  };
  const handleCreateClick = () => {
    setNoteToEdit(null);
    openModal();
  };
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this note?')) {
      try {
        await deleteNote(id);
        fetchNotes();
      } catch (err) {
        setError("Failed to delete note.");
      }
    }
  };

  if (loading) return <Loader />;
  if (error) return <Alert color="red" title="Error" icon={<IconAlertCircle />}>{error}</Alert>;

  return (
    <>
      <CreateNoteModal opened={modalOpened} onClose={closeModal} onSuccess={fetchNotes} noteToEdit={noteToEdit} />
      <Group justify="space-between" mb="md">
        <Title order={2}>My Notes</Title>
        <Button onClick={handleCreateClick}>Add New Note</Button>
      </Group>

      <SimpleGrid cols={{ base: 1, sm: 2, lg: 3 }}>
        {notes.length > 0 ? (
          notes.map((note) => (
            <Card shadow="sm" padding="lg" radius="md" withBorder key={note.id}>
              <Text size="sm" c="dimmed">{new Date(note.note_date).toLocaleString()}</Text>
              <Text mt="xs" mb="md">{note.note_text}</Text>
              <Group justify="flex-end">
                <ActionIcon variant="subtle" onClick={() => handleEditClick(note)}><IconPencil size={18} /></ActionIcon>
                <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteClick(note.id)}><IconTrash size={18} /></ActionIcon>
              </Group>
            </Card>
          ))
        ) : (
          <Text>No notes found. Add one!</Text>
        )}
      </SimpleGrid>
    </>
  );
};

export default NotesPage;