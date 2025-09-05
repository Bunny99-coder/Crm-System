// Replace the contents of CreateLeadModal.tsx

import { Modal, Button, Select, Stack, Textarea } from '@mantine/core';
import { useForm, isNotEmpty } from '@mantine/form';
import { useState, useEffect } from 'react';
import { createLead, updateLead, getContacts, getUsers, getProperties } from '../services/apiService';
import type { Lead, CreateLeadPayload, Contact, UserSelectItem, Property } from '../services/apiService';

interface CreateLeadModalProps {
  opened: boolean;
  onClose: () => void;
  onSuccess: () => void;
  leadToEdit: Lead | null; // <-- NEW
}

const CreateLeadModal = ({ opened, onClose, onSuccess, leadToEdit }: CreateLeadModalProps) => {
  const isEditing = !!leadToEdit;
  const [contacts, setContacts] = useState<{ value: string; label: string }[]>([]);
  const [users, setUsers] = useState<{ value: string; label: string }[]>([]);
  const [properties, setProperties] = useState<{ value: string; label: string }[]>([]);
  const [loadingData, setLoadingData] = useState(true);

  // Hardcoded sources and statuses
  const leadSources = [
    { value: '1', label: 'Website Inquiry' }, { value: '2', label: 'Phone Call' },
    { value: '3', label: 'Social Media' }, { value: '4', label: 'Referral' },
  ];
  const leadStatuses = [
    { value: '1', label: 'New' }, { value: '2', label: 'Contacted' },
    { value: '3', label: 'Qualified' }, { value: '4', label: 'Converted' }, { value: '5', label: 'Lost' },
  ];

  const form = useForm({
    initialValues: {
      contact_id: '',
      property_id: null as string | null, // Important for clearable Select
      source_id: '1',
      status_id: '1',
      assigned_to: '',
      notes: '',
    },
    validate: {
      contact_id: isNotEmpty('A contact must be selected'),
      assigned_to: isNotEmpty('A user must be assigned'),
    },
  });

  useEffect(() => {
    if (opened) {
      const fetchData = async () => { /* ... (fetchData is the same as before) ... */ };
      fetchData();

      // Pre-fill form if in edit mode
      if (isEditing) {
        form.setValues({
          contact_id: String(leadToEdit.contact_id),
          property_id: leadToEdit.property_id ? String(leadToEdit.property_id) : null,
          source_id: String(leadToEdit.source_id),
          status_id: String(leadToEdit.status_id),
          assigned_to: String(leadToEdit.assigned_to),
          notes: leadToEdit.notes || '',
        });
      } else {
        form.reset();
      }
    }
  }, [leadToEdit, opened]);

  // Fetch data for dropdowns (same as before)
  useEffect(() => {
    if (opened) {
      const fetchData = async () => {
        setLoadingData(true);
        try {
          const [contactsData, usersData, propertiesData] = await Promise.all([ getContacts(), getUsers(), getProperties() ]);
          setContacts(contactsData.map((c: Contact) => ({ value: String(c.id), label: `${c.first_name} ${c.last_name}` })));
          setUsers(usersData.map((u: UserSelectItem) => ({ value: String(u.id), label: u.username })));
          setProperties(propertiesData.map((p: Property) => ({ value: String(p.id), label: p.name })));
        } catch (error) { console.error("Failed to fetch data for lead form", error); }
        finally { setLoadingData(false); }
      };
      fetchData();
    }
  }, [opened]);

  const handleSubmit = async (values: typeof form.values) => {
    try {
      const payload: CreateLeadPayload = { // CreateLeadPayload and UpdateLeadPayload are the same shape
        contact_id: Number(values.contact_id),
        property_id: values.property_id ? Number(values.property_id) : undefined,
        source_id: Number(values.source_id),
        status_id: Number(values.status_id),
        assigned_to: Number(values.assigned_to),
        notes: values.notes,
      };
      
      if (isEditing) {
        await updateLead(leadToEdit.id, payload);
      } else {
        await createLead(payload);
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Failed to save lead", error);
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title={isEditing ? 'Edit Lead' : 'Create New Lead'} centered>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <Select label="Contact" data={contacts} searchable required disabled={loadingData} {...form.getInputProps('contact_id')} />
          <Select label="Assigned To" data={users} searchable required disabled={loadingData} {...form.getInputProps('assigned_to')} />
          <Select label="Property (Optional)" data={properties} searchable clearable disabled={loadingData} {...form.getInputProps('property_id')} />
          <Select label="Lead Source" data={leadSources} required {...form.getInputProps('source_id')} />
          <Select label="Lead Status" data={leadStatuses} required {...form.getInputProps('status_id')} />
          <Textarea label="Notes" {...form.getInputProps('notes')} />
          <Button type="submit" mt="md" loading={loadingData}>{isEditing ? 'Update Lead' : 'Create Lead'}</Button>
        </Stack>
      </form>
    </Modal>
  );
};

export default CreateLeadModal;