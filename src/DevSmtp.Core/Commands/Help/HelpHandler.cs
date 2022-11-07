using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class HelpHandler : ICommandHandler<Help, HelpResult>
    {
        private readonly IDataStore _dataStore;

        public HelpHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<HelpResult> ExecuteAsync(Help command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
