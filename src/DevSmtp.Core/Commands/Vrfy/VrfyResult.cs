namespace DevSmtp.Core.Commands
{
    public sealed class VrfyResult : CommandResult
    {
        public VrfyResult()
        {
        }

        public VrfyResult(Exception error)
            : base(error)
        {
        }
    }
}
