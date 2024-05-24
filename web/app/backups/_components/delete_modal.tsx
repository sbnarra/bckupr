import {Backup, ModelError as Error, DeleteBackupRequest} from '../../../components/spec'
import {NewBackupApi} from '../../../components/api'
import React from "react";
import {Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure} from "@nextui-org/react";

export default function DeleteBackup(props: {
  backup: Backup
  onDelete: () => void
}) {
  const {isOpen, onOpen, onOpenChange} = useDisclosure();
  const deleteBackup = () => {
    NewBackupApi().deleteBackup({
      id: props.backup.id
    }).then(() => {
      alert("Deleted " + props.backup.id)
      props.onDelete()
    }).catch((err) => {
      alert("Error: " + err.error)
    })
  }

  return (
      <Button onPress={onOpen}>
        Delete
        
      <Modal 
        isOpen={isOpen} 
        placement="auto"
        onOpenChange={onOpenChange} 
      >
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1">Delete Backup</ModalHeader>
              <ModalBody>
                <p> 
                  Danger! Danger! Danger!
                </p>
                <p> 
                  This can&apos;t be undone...
                </p>
              </ModalBody>
              <ModalFooter>
                <Button color="primary" onPress={onClose}>
                  Cancel
                </Button>
                <Button color="danger" variant="light" onPress={function() {
                  deleteBackup()
                  onClose()
                }}>Delete</Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
      </Button>
    );
}
