package proxmox

import (
	"testing"

	proxmoxv1alpha1 "github.com/alperencelik/kubemox/api/proxmox/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckVMType_Template(t *testing.T) {
	vm := &proxmoxv1alpha1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-template",
			Namespace: "default",
		},
		Spec: proxmoxv1alpha1.VirtualMachineSpec{
			Template: &proxmoxv1alpha1.VirtualMachineSpecTemplate{
				Name: "test-template",
			},
		},
	}

	vmType := CheckVMType(vm)

	assert.Equal(t, VirtualMachineTemplateType, vmType)
}

func TestCheckVMType_Undefined(t *testing.T) {
	vm := &proxmoxv1alpha1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-vm",
			Namespace: "default",
		},
		Spec: proxmoxv1alpha1.VirtualMachineSpec{},
	}

	vmType := CheckVMType(vm)

	assert.Equal(t, "undefined", vmType)
}

func TestSortDisks_SingleDisk(t *testing.T) {
	disks := []proxmoxv1alpha1.VirtualMachineDisk{
		{
			Device:  "scsi0",
			Storage: "local-lvm",
			Size:    32,
		},
	}

	sorted := sortDisks(disks)

	assert.Len(t, sorted, 1)
	assert.Equal(t, "scsi0", sorted[0].Device)
}

func TestSortDisks_MultipleDisks_SortByStorage(t *testing.T) {
	disks := []proxmoxv1alpha1.VirtualMachineDisk{
		{
			Device:  "scsi1",
			Storage: "ceph-pool",
			Size:    64,
		},
		{
			Device:  "scsi0",
			Storage: "local-lvm",
			Size:    32,
		},
	}

	sorted := sortDisks(disks)

	assert.Len(t, sorted, 2)
	assert.Equal(t, "ceph-pool", sorted[0].Storage)
	assert.Equal(t, "local-lvm", sorted[1].Storage)
}

func TestSortDisks_MultipleDisks_SortByDevice(t *testing.T) {
	disks := []proxmoxv1alpha1.VirtualMachineDisk{
		{
			Device:  "scsi2",
			Storage: "local-lvm",
			Size:    32,
		},
		{
			Device:  "scsi0",
			Storage: "local-lvm",
			Size:    64,
		},
		{
			Device:  "scsi1",
			Storage: "local-lvm",
			Size:    128,
		},
	}

	sorted := sortDisks(disks)

	assert.Len(t, sorted, 3)
	assert.Equal(t, "scsi0", sorted[0].Device)
	assert.Equal(t, "scsi1", sorted[1].Device)
	assert.Equal(t, "scsi2", sorted[2].Device)
}

func TestSortDisks_EmptyList(t *testing.T) {
	disks := []proxmoxv1alpha1.VirtualMachineDisk{}

	sorted := sortDisks(disks)

	assert.Empty(t, sorted)
}

func TestSortDisks_NilList(t *testing.T) {
	var disks []proxmoxv1alpha1.VirtualMachineDisk

	sorted := sortDisks(disks)

	assert.Empty(t, sorted)
}

func TestVirtualMachineConstants(t *testing.T) {
	assert.Equal(t, "running", VirtualMachineRunningState)
	assert.Equal(t, "stopped", VirtualMachineStoppedState)
	assert.Equal(t, "template", VirtualMachineTemplateType)
	assert.Equal(t, "scratch", VirtualMachineScratchType)
	assert.Equal(t, 10, AgentTimeoutSeconds)
	assert.Equal(t, 15, virtualMachineCreateTimesNum)
	assert.Equal(t, 10, virtualMachineStartTimesNum)
	assert.Equal(t, 10, virtualMachineStopTimesNum)
	assert.Equal(t, 10, virtualMachineDeleteTimesNum)
}

func TestVirtualMachineTags(t *testing.T) {
	assert.NotEmpty(t, virtualMachineTag)
	assert.NotEmpty(t, ManagedVirtualMachineTag)
	assert.NotEmpty(t, virtualMachineTemplateTag)

	assert.Equal(t, "kubemox", virtualMachineTag)
	assert.Equal(t, "kubemox-managed-vm", ManagedVirtualMachineTag)
	assert.Equal(t, "kubemox-template", virtualMachineTemplateTag)
}

func TestCreateVMSemaphore(t *testing.T) {
	assert.NotNil(t, createVMSemaphore)
	assert.Equal(t, 10, cap(createVMSemaphore), "Semaphore should allow 10 concurrent VM creations")
}

func TestEnableProxmoxTaskLogs(t *testing.T) {
	assert.False(t, EnableProxmoxTaskLogs, "Task logs should be disabled by default")
}

func TestVirtualMachineDeltaComparison(t *testing.T) {
	vm1 := VirtualMachineComparison{
		Cores:   4,
		Sockets: 2,
		Memory:  8192,
	}

	vm2 := VirtualMachineComparison{
		Cores:   4,
		Sockets: 2,
		Memory:  8192,
	}

	vm3 := VirtualMachineComparison{
		Cores:   8,
		Sockets: 2,
		Memory:  8192,
	}

	assert.Equal(t, vm1.Cores, vm2.Cores)
	assert.Equal(t, vm1.Sockets, vm2.Sockets)
	assert.Equal(t, vm1.Memory, vm2.Memory)
	assert.NotEqual(t, vm1.Cores, vm3.Cores)
}
