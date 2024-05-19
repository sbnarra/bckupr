import { useEffect, useMemo, useState } from "react";
import {
  Selection,
  Listbox,
  Checkbox,
  ListboxItem,
  ListboxSection, 
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
import { Backup, TaskInput, Filters, Status, Volume } from '@/components/spec';

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
  
  var filters = new Filters([], [], [], [])
  var taskInput = new TaskInput(filters)

  const startBackup = (onClose: () => void) => {
    console.log(JSON.stringify(taskInput))
    var api = NewBackupApi()
    api.startBackup(taskInput, (err: any, backup: Backup) => {
      if (err === null) {
        props.setId(backup.id)
        onClose()
        props.progressOnOpen()
      } else {
        alert("Error: " + err.error)
      }
    })
  }

  const [selectedKeys, setSelectedKeys] = useState(new Set(["labelled", "writers", "linked"]));

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
              <Checkbox radius="full" onChange={(e) => taskInput.no_dry_run = e.target.checked}>
                No Dry Run
              </Checkbox>
              
              <ListInput label="Include Names" updated={(s) => filters.include_names = s}/>
              <ListInput label="Include Volumes" updated={(s) => filters.include_volumes = s}/>
              <ListInput label="Exclude Names" updated={(s) => filters.exclude_names = s}/>
              <ListInput label="Exclude Volumes" updated={(s) => filters.exclude_volumes = s}/>

              <Listbox 
                aria-label="Stop Modes"
                disallowEmptySelection
                selectionMode="multiple"
                selectedKeys={selectedKeys}
                onSelectionChange={(keys: Selection) => {
                  setSelectedKeys(keys as Set<string>)
                  taskInput.stop_modes = keys
                }}
              >
                <ListboxSection title="Stop Modes"> 
                  <ListboxItem key="all">All</ListboxItem>
                  <ListboxItem key="labelled">Labelled</ListboxItem>
                  <ListboxItem key="writers">Writers</ListboxItem>
                  <ListboxItem key="attached">Attached</ListboxItem>
                  <ListboxItem key="linked">Linked</ListboxItem>
                </ListboxSection>
              </Listbox>

              <Input
                label="Label Prefix"
                defaultValue="bckupr"
                onChange={(e) => taskInput.label_prefix = e.target.value}
              />
            </ModalBody>
            <ModalFooter>
              <Button color="danger" variant="light" onPress={onClose}>
                Cancel
              </Button>
              <Button color="primary" onPress={() => startBackup(onClose)}>
                Start
              </Button>
              
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>)
}

function ListInput(props: {
  label: string
  updated: (items: string[]) => void
}) {
  const [inputted, setInputted] = useState("")
  const [inputs, setInputs] = useState(new Set<string>())

  return (<>
    <div className="flex">
      <Input isClearable
        label={props.label}
        placeholder="..."
        value={inputted}
        onClear={() => setInputted("")}
        onChange={(e) => setInputted(e.target.value)}
      />
      <div className="content-around">
        <Button className="h-full min-w" onClick={() => {
          inputs.add(inputted)
          setInputted("")
          const updated = new Set(inputs)
          setInputs(updated)
          props.updated(Array.from(updated))
        }}>+</Button>
      </div>
    </div>
    {Array.from(inputs).map((name: string) => <div key={name} className="flex">
      <Input defaultValue={name} isReadOnly/>
      <div className="content-around">
        <Button className="h-full min-w" onClick={() => {
          inputs.delete(name)
          const updated = new Set(inputs)
          setInputs(updated)
          props.updated(Array.from(updated))
        }}>-</Button>
      </div>
    </div>)}
  </>)
}

function BackupWaitModal(props: {
  id: string
  onOpenChange: () => void
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
              <Button color="danger" variant="light" onPress={onClose}>
                Cancel
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>)
}