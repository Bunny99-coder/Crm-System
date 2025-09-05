import { Modal, Button, TextInput, NumberInput, Select, Stack } from '@mantine/core';
import { useForm, isNotEmpty } from '@mantine/form';
import { useEffect } from 'react';
import { createProperty, updateProperty } from '../services/apiService';
import type { Property, CreatePropertyPayload } from '../services/apiService';

interface CreatePropertyModalProps {
  opened: boolean;
  onClose: () => void;
  onSuccess: () => void;
  propertyToEdit: Property | null;
}

const CreatePropertyModal = ({ opened, onClose, onSuccess, propertyToEdit }: CreatePropertyModalProps) => {
  const isEditing = !!propertyToEdit;

  const sites = [{ value: '1', label: 'Downtown Condos' }];
  const propertyTypes = [
    { value: '1', label: 'Apartment' }, { value: '2', label: 'Villa' },
    { value: '3', label: 'Office' }, { value: '4', label: 'Townhouse' },
  ];
  const statuses = ['Available', 'Sold', 'Pending'];

  // --- THE FIX IS HERE (Part 1) ---
  // We will manage the form state for Selects as strings.
  const form = useForm({
    initialValues: {
      name: '',
      site_id: '1', // The Select component works best with string values
      property_type_id: '1',
      unit_no: '',
      price: 0,
      status: 'Available',
    },
    validate: {
      name: isNotEmpty('Property name is required'),
      price: (value) => (value > 0 ? null : 'Price must be greater than zero'),
    },
  });

  useEffect(() => {
    if (isEditing) {
      form.setValues({
        name: propertyToEdit.name,
        // Convert the number from the API to a string for the form
        site_id: String(propertyToEdit.site_id),
        property_type_id: String(propertyToEdit.property_type_id),
        unit_no: propertyToEdit.unit_no || '',
        price: propertyToEdit.price,
        status: propertyToEdit.status,
      });
    } else {
      form.reset();
    }
  }, [propertyToEdit, opened]);

  const handleSubmit = async (values: typeof form.values) => {
    try {
      // --- THE FIX IS HERE (Part 2) ---
      // Convert the string values from the form back to numbers for the API payload.
      const payload: CreatePropertyPayload = {
        name: values.name,
        site_id: Number(values.site_id),
        property_type_id: Number(values.property_type_id),
        unit_no: values.unit_no,
        price: values.price,
        status: values.status,
      };

      if (isEditing) {
        await updateProperty(propertyToEdit.id, payload);
      } else {
        await createProperty(payload);
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Failed to save property", error);
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title={isEditing ? 'Edit Property' : 'Create New Property'} centered>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput label="Property Name" required {...form.getInputProps('name')} />

          {/* --- THE FIX IS HERE (Part 3) --- */}
          {/* We remove the problematic options object from getInputProps for the Selects */}
          <Select label="Site" data={sites} required {...form.getInputProps('site_id')} />
          <Select label="Property Type" data={propertyTypes} required {...form.getInputProps('property_type_id')} />
          
          <TextInput label="Unit No." {...form.getInputProps('unit_no')} />
          <NumberInput label="Price" prefix="$ " required thousandSeparator min={0} {...form.getInputProps('price')} />
          <Select label="Status" data={statuses} required {...form.getInputProps('status')} />
          <Button type="submit" mt="md">{isEditing ? 'Update Property' : 'Create Property'}</Button>
        </Stack>
      </form>
    </Modal>
  );
};

export default CreatePropertyModal;