import React, { useEffect, useState } from "react";
import {Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure} from "@nextui-org/react";
import { NewBackupApi } from '@/components/api';
import { Backup, TaskInput, Filters } from '@/components/spec';

export function CreateBackup(props: {}) {
  const {isOpen: optionsIsOpen, onOpen: optionsOnOpen, onOpenChange: optionsOnOpenChange} = useDisclosure();
  const {isOpen: progressIsOpen, onOpen: progressOnOpen, onOpenChange: progressOnOpenChange} = useDisclosure();
  const [id, setId] = useState("")

  return (<div className="flex flex-col gap-2">
      <Button onPress={optionsOnOpen}>Create Backup</Button>
      {optionsIsOpen && <BackupOptionsModal
        progressOnOpen={progressOnOpen}
        onOpenChange={optionsOnOpenChange}
        setId={setId}
      />}
      <p>{optionsIsOpen}</p>
      <p>{progressIsOpen}</p>
      {progressIsOpen && <BackupWaitModal
        id={id}
        onOpenChange={progressOnOpenChange}
      />}
    </div>);
}

function BackupOptionsModal(props: {
  progressOnOpen: () => void
  onOpenChange: () => void
  setId: (id: string) => void
}) {
  
  var f = new Filters([], [], [], [])
  var t = new TaskInput(f)

  function startBackup(onClose: () => void) {
    var api = NewBackupApi()
    api.startBackup(t, function(err: any, backup: Backup) {
      if (err === null) {
        props.setId(backup.id)
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
            <ModalHeader className="flex flex-col gap-1">Create Backup</ModalHeader>
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
                startBackup(onClose)
              }}>
                Action
              </Button>
              
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>)
}

function BackupWaitModal(props: {
  id: string
  onOpenChange: () => void
}) {

  const [backup, setBackup] = useState<Backup>()
  const [error, setError] = useState<String>()

  useEffect(() => {
    const id = setInterval(() => {
      const api = NewBackupApi()
      api.getBackup(props.id, function(err: any, backup: Backup) {
        if (err != null) {
          setError(err.response.text)
        } else {
          setBackup(backup)
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
            <ModalHeader className="flex flex-col gap-1">Waiting for Backup</ModalHeader>
            <ModalBody>
              {error && <pre>{error}</pre>}
              {backup && <pre>{JSON.stringify(backup)}</pre>}
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