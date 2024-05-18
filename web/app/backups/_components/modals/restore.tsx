import React, { useEffect, useState } from "react";
import {Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure} from "@nextui-org/react";
import { NewBackupApi } from '@/components/api';
import { Backup, TaskInput, Filters, Error, Restore } from '@/components/spec';

export default function RestoreBackup(props: {
  backup: Backup
}) {
  const {isOpen: optionsIsOpen, onOpen: optionsOnOpen, onOpenChange: optionsOnOpenChange} = useDisclosure();
  const {isOpen: progressIsOpen, onOpen: progressOnOpen, onOpenChange: progressOnOpenChange} = useDisclosure();

  return (<div className="flex flex-col gap-2">
      <Button onPress={optionsOnOpen}>Restore</Button>
      {optionsIsOpen && <RestoreOptionsModal
        backup={props.backup}
        progressOnOpen={progressOnOpen}
        onOpenChange={optionsOnOpenChange}
      />}
      <p>{optionsIsOpen}</p>
      <p>{progressIsOpen}</p>
      {progressIsOpen && <RestoreWaitModal
        backup={props.backup}
        onOpenChange={progressOnOpenChange}
      />}
    </div>);
}

function RestoreOptionsModal(props: {
  backup: Backup
  progressOnOpen: () => void
  onOpenChange: () => void
}) {
  
  var f = new Filters([], [], [], [])
  var t = new TaskInput(f)

  function startRestore(onClose: () => void) {
    var api = NewBackupApi()
    api.startRestore(props.backup.id, t, function(err: Error) {
      if (err === null) {
        onClose()
        props.progressOnOpen()
      } else {
        alert("Error: " + err.error)
      }
    })
  }

  return (<Modal 
      isOpen={true} 
      placement="auto"
      onOpenChange={props.onOpenChange} 
    >
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className="flex flex-col gap-1">Restore Backup</ModalHeader>
            <ModalBody>
              <p>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit.
                Nullam pulvinar risus non risus hendrerit venenatis.
                Pellentesque sit amet hendrerit risus, sed porttitor quam.
              </p>
              <p>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit.
                Nullam pulvinar risus non risus hendrerit venenatis.
                Pellentesque sit amet hendrerit risus, sed porttitor quam.
              </p>
            </ModalBody>
            <ModalFooter>
              <Button color="danger" variant="light" onPress={onClose}>
                Close
              </Button>
              <Button color="primary" onPress={function() {
                startRestore(onClose)
              }}>
                Action
              </Button>
              
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>)
}

function RestoreWaitModal(props: {
  backup: Backup
  onOpenChange: () => void
}) {

  const [restore, setRestore] = useState<Restore>()
  const [error, setError] = useState<String>()

  useEffect(() => {
    const id = setInterval(() => {
      const api = NewBackupApi()
      api.getRestore(props.backup.id, function(err: any, restore: Restore) {
        if (err != null) {
          setError(err.response.text)
        } else {
          setRestore(restore)
          setError(undefined)
        }
      })
    }, 1000)
    return () => clearInterval(id)
  }, [])

  return (<Modal 
      isOpen={true} 
      placement="auto"
      onOpenChange={props.onOpenChange} 
    >
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className="flex flex-col gap-1">Waiting for Restore</ModalHeader>
            <ModalBody>
              {error && <pre>{error}</pre>}
              {restore && <pre>{JSON.stringify(restore)}</pre>}
            </ModalBody>
            <ModalFooter>
              <Button color="danger" variant="light" onPress={function() {
                onClose()
              }}>
                Cancel
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>)
}