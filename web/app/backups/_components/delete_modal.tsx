import {Backup, Error} from '../../../components/spec'
import {NewBackupApi} from '../../../components/api'
import React from "react";
import {Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure} from "@nextui-org/react";

export default function DeleteBackup(props: {
  backup: Backup
}) {
  const {isOpen, onOpen, onOpenChange} = useDisclosure();

  function deleteBackup() {
    NewBackupApi().deleteBackup(props.backup.id, function(err: Error) {
      if (err === null) {
        alert("Deleted " + props.backup.id)
      } else {
        alert("Error: " + err.error)
      }
    })
  }

  return (
    <div className="flex flex-col gap-2">
      <Button onPress={onOpen}>Delete</Button>
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
                <Button color="danger" variant="light" onPress={onClose}>
                  Cancel
                </Button>
                <Button color="primary"  onPress={function() {
                  deleteBackup()
                  onClose()
                }}></Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </div>
  );
}
