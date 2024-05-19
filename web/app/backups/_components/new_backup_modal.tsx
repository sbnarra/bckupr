import { useEffect, useState } from "react";
import {
  Modal, 
  ModalContent, 
  ModalHeader, 
  ModalBody, 
  ModalFooter, 
  Button, 
  useDisclosure, 
  Input
} from "@nextui-org/react";
import { NewBackupApi } from '@/components/api';
import { Backup, Status, Volume, TaskInput } from '@/components/spec';
import { TaskInputModal } from './task_input_modal'

export function CreateBackup(props: {
  refresh: () => void
}) {
  const {isOpen: optionsIsOpen, onOpen: optionsOnOpen, onOpenChange: optionsOnOpenChange} = useDisclosure();
  const {isOpen: progressIsOpen, onOpen: progressOnOpen, onOpenChange: progressOnOpenChange} = useDisclosure();
  const [id, setId] = useState("")

  return (<div className="flex flex-col gap-2">
      <Button onPress={optionsOnOpen}>Create Backup</Button>
      {optionsIsOpen && <TaskInputModal
        progressOnOpen={progressOnOpen}
        onOpenChange={optionsOnOpenChange}
        callApi={(taskInput: TaskInput, onClose: () => void) => {
          console.log(JSON.stringify(taskInput))
          var api = NewBackupApi()
          api.startBackup(taskInput, (err: any, backup: Backup) => {
            if (err === null) {
              setId(backup.id)
              onClose()
              progressOnOpen()
            } else {
              alert("Error: " + err.error)
            }
          })}}
      />}
      <p>{optionsIsOpen}</p>
      <p>{progressIsOpen}</p>
      {progressIsOpen && <BackupWaitModal
        id={id}
        refresh={props.refresh}
        onOpenChange={progressOnOpenChange}
      />}
    </div>);
}

function BackupWaitModal(props: {
  id: string
  onOpenChange: () => void
  refresh: () => void
}) {

  const [backup, setBackup] = useState<Backup>()
  const [error, setError] = useState<String>()

  useEffect(() => {
    const api = NewBackupApi()

    const id = setInterval(() => {
      api.getBackup(props.id, (err: any, backup: Backup) => {
        if (err != null) {
          setError(err.response.text)
        } else {
          setBackup(backup)
          setError(undefined)

          const s = new Status()
          if (backup.status == s.error || backup.status == s.completed) {
            clearInterval(id)
          }
        }
      })
    }, 1000)
    return () => clearInterval(id)
  }, [props.id])

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
              {error && <pre>Error: {error}</pre>}
              {backup && <>
              <p>Id: {backup.id}</p>
              <p>Created: {JSON.stringify(backup.created)}</p>
              <p>Type: {backup.type}</p>
              <p>Status: {backup.status}</p>
              {backup.error && <p>Error: {backup.error}</p>}
              {backup.volumes && <p>Volumes:
                <ul>
                  {backup.volumes.map((volume: Volume) => <li key={volume.name}>{JSON.stringify(volume)}</li>)}
                </ul>
              </p>}
              </>}
            </ModalBody>
            <ModalFooter>
              <Button color="danger" variant="light" onPress={() => {
                onClose()
                props.refresh()
              }}>
                Cancel
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>)
}